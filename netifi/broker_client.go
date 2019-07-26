package netifi

import (
	"encoding/base64"
	"github.com/google/uuid"
	"github.com/netifi/netifi-go/discovery_strategy"
	"github.com/netifi/netifi-go/internal/framing"
	"github.com/netifi/netifi-go/internal/rsocket/named_rsocket_wrapper"
	"github.com/netifi/netifi-go/internal/rsocket/reconnecting_rsocket"
	"github.com/netifi/netifi-go/internal/rsocket/transforming_rsocket"
	"github.com/netifi/netifi-go/tags"
	"github.com/pkg/errors"
	"github.com/rsocket/rsocket-go"
	"github.com/rsocket/rsocket-go/payload"
	rrpc "github.com/rsocket/rsocket-rpc-go"
	"math/rand"
	"net"
	"sync"
)

type Builder interface {
	DiscoveryStrategy(strategy discovery_strategy.DiscoveryStrategy) Builder
	AccessKey(accessKey uint64) Builder
	AccessToken(accessToken []byte) Builder
	AccessTokenBase64(accessToken string) Builder
	Uri(uri string) Builder
	Group(g string) Builder
	Tags(tags tags.Tags) Builder
	Destination(destination string) Builder
	Public(p bool) Builder
	IP(ip net.IP) Builder
	Build() (BrokerClient, error)
}

func New() Builder {
	return &brokerClientConfig{}
}

type brokerClientConfig struct {
	rh           rrpc.RequestHandlingRSocket
	dest         string
	ds           discovery_strategy.DiscoveryStrategy
	group        string
	host         string
	port         int
	uri          string
	key          uint64
	token        []byte
	encodedToken string
	tags         tags.Tags
	ip           net.IP
	flags        uint32
	uuid         uuid.UUID
}

func (b *brokerClientConfig) Public(p bool) Builder {
	if p {
		b.flags = 1
	} else {
		b.flags = 0
	}
	return b
}

func (b *brokerClientConfig) Group(g string) Builder {
	b.group = g
	return b
}

func (b *brokerClientConfig) IP(ip net.IP) Builder {
	b.ip = ip
	return b
}

func (b *brokerClientConfig) Uri(uri string) Builder {
	b.uri = uri
	return b
}
func (b *brokerClientConfig) AccessKey(accessKey uint64) Builder {
	b.key = accessKey
	return b
}

func (b *brokerClientConfig) AccessToken(accessToken []byte) Builder {
	b.token = accessToken
	return b
}

func (b *brokerClientConfig) AccessTokenBase64(accessToken string) Builder {
	b.encodedToken = accessToken
	return b
}

func (b *brokerClientConfig) DiscoveryStrategy(strategy discovery_strategy.DiscoveryStrategy) Builder {
	b.ds = strategy
	return b
}

func (b *brokerClientConfig) Tags(tags tags.Tags) Builder {
	b.tags = tags
	return b
}

func (b *brokerClientConfig) Destination(destination string) Builder {
	b.dest = destination
	return b
}

func (b *brokerClientConfig) Build() (client BrokerClient, e error) {
	var selector RSocketSelector

	b.uuid = uuid.New()

	if b.ds != nil {
		selector = &simpleRSocketSelector{config: *b}
	} else if len(b.uri) > 0 {
		b.ds = discovery_strategy.New(b.uri)
		selector = &simpleRSocketSelector{config: *b}
	} else {
		e = errors.New("must include either uri or discovery service")
		return
	}

	if b.key < 1 {
		e = errors.New("access key must be greater than 0")
		return
	}

	if len(b.encodedToken) < 1 {
		if len(b.token) < 1 {
			e = errors.New("must include access token")
			return
		}
	} else {
		var bytes []byte
		bytes, e = base64.StdEncoding.DecodeString(b.encodedToken)
		if e != nil {
			return
		}
		b.token = bytes
	}

	if len(b.group) < 1 {
		e = errors.New("must include a group")
		return
	}

	var destinationTag *tags.Tag
	if len(b.dest) > 0 {
		destinationTag = tags.New("com.netifi.destination", b.dest)
	} else {
		destinationTag = tags.New("com.netifi.destination", b.uuid.String())
	}

	b.tags = b.tags.And(destinationTag)

	client = &brokerClient{
		selector: selector,
		config:   *b,
	}
	return
}

type BrokerSocket rsocket.RSocket

type BrokerClient interface {
	AddService(rs rrpc.RrpcRSocket) error

	AddNamedSocket(name string, rs rsocket.RSocket) error

	GroupServiceSocket(group string, tags tags.Tags) BrokerSocket

	BroadcastServiceSocket(group string, tags tags.Tags) BrokerSocket

	ShardServiceSocket(group string, shardKey []byte, tags tags.Tags) BrokerSocket

	GroupNamedRSocket(group string, tags tags.Tags) BrokerSocket

	BroadcastNamedRSocket(group string, tags tags.Tags) BrokerSocket

	ShardNamedRSocket(group string, shardKey []byte, tags tags.Tags) BrokerSocket
}

type RSocketSelector interface {
	selectRSocket() rsocket.RSocket
}

type simpleRSocketSelector struct {
	RSocketSelector
	sync.Mutex
	config brokerClientConfig
	rs     rsocket.RSocket
}

func (s *simpleRSocketSelector) selectRSocket() rsocket.RSocket {
	s.Lock()
	defer s.Unlock()

	if s.rs != nil {
		return s.rs
	} else {
		return reconnecting_rsocket.New(func() string {
			nodes := <-s.config.ds.DiscoverNodes()
			l := len(nodes)
			i := rand.Intn(l)
			return nodes[i]
		}, s.config.rh, func() payload.Payload {
			d, e := framing.EncodeDestinationSetup(s.config.ip,
				s.config.group,
				s.config.key,
				s.config.token,
				s.config.uuid,
				s.config.flags, s.config.tags)
			if e != nil {
				panic("error encoding destination setup")
			}
			return payload.New(nil, d)
		})
	}
}

type brokerClient struct {
	BrokerClient
	config   brokerClientConfig
	selector RSocketSelector
}

func (b *brokerClient) AddService(rs rrpc.RrpcRSocket) error {
	return b.config.rh.Register(rs)
}

func (b *brokerClient) AddNamedSocket(name string, rs rsocket.RSocket) error {
	wrapper := named_rsocket_wrapper.New(name, rs)
	return b.config.rh.Register(wrapper)
}

func (b *brokerClient) GroupServiceSocket(group string, tags tags.Tags) BrokerSocket {
	socket := &transforming_rsocket.PayloadTransformingRSocket{
		Transformer: func(msg payload.Payload) (p payload.Payload, e error) {
			data := msg.Data()
			metadata, _ := msg.Metadata()
			md, e := framing.EncodeGroup(group, metadata, tags)
			if e != nil {
				return
			}
			p = payload.New(data, md)
			return
		},
		Source: b.selector.selectRSocket,
	}

	return socket
}

func (b *brokerClient) BroadcastServiceSocket(group string, tags tags.Tags) BrokerSocket {
	socket := &transforming_rsocket.PayloadTransformingRSocket{
		Transformer: func(msg payload.Payload) (p payload.Payload, e error) {
			data := msg.Data()
			metadata, _ := msg.Metadata()
			md, e := framing.EncodeBroadcast(group, metadata, tags)
			if e != nil {
				return
			}
			p = payload.New(data, md)
			return
		},
		Source: b.selector.selectRSocket,
	}

	return socket
}

func (b *brokerClient) ShardServiceSocket(group string, shardKey []byte, tags tags.Tags) BrokerSocket {
	socket := &transforming_rsocket.PayloadTransformingRSocket{
		Transformer: func(msg payload.Payload) (p payload.Payload, e error) {
			data := msg.Data()
			metadata, _ := msg.Metadata()
			md, e := framing.EncodeShard(group, metadata, shardKey, tags)
			if e != nil {
				return
			}
			p = payload.New(data, md)
			return
		},
		Source: b.selector.selectRSocket,
	}

	return socket
}

func (b *brokerClient) GroupNamedRSocket(group string, tags tags.Tags) BrokerSocket {
	return named_rsocket_wrapper.New(group, b.GroupServiceSocket(group, tags))
}

func (b *brokerClient) BroadcastNamedRSocket(group string, tags tags.Tags) BrokerSocket {
	return named_rsocket_wrapper.New(group, b.BroadcastServiceSocket(group, tags))
}

func (b *brokerClient) ShardNamedRSocket(group string, shardKey []byte, tags tags.Tags) BrokerSocket {
	return named_rsocket_wrapper.New(group, b.ShardServiceSocket(group, shardKey, tags))
}