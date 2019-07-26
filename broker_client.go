package netifi_go

import (
	"github.com/netifi/netifi-go/internal/framing"
	rsocket2 "github.com/netifi/netifi-go/internal/rsocket"
	"github.com/rsocket/rsocket-go"
	"github.com/rsocket/rsocket-go/payload"
	rrpc "github.com/rsocket/rsocket-rpc-go"
	"math/rand"
	"sync"
)

type Builder interface {
	DiscoveryStrategy(strategy DiscoveryStrategy) Builder
	AccessKey(accessKey string) Builder
	AccessToken(accessToken []byte) Builder
	Uri(uri string) Builder
	Tags(tags Tags) Builder
	Destination(destination string) Builder
	Build() BrokerClient
}

type brokerClientConfig struct {
	rh    rrpc.RequestHandlingRSocket
	ds    DiscoveryStrategy
	host  string
	port  int
	uri   string
	key   string
	token []byte
	tags  Tags
}

func (b *brokerClientConfig) Uri(uri string) Builder {
	b.uri = uri
	return b
}
func (b *brokerClientConfig) AccessKey(accessKey string) Builder {
	b.key = accessKey
	return b
}

func (b *brokerClientConfig) AccessToken(accessToken []byte) Builder {
	b.token = accessToken
	return b
}

func (b *brokerClientConfig) DiscoveryStrategy(strategy DiscoveryStrategy) Builder {
	b.ds = strategy
	return b
}

func (b *brokerClientConfig) Tags(tags Tags) Builder {
	return b
}

func (b *brokerClientConfig) Destination(destination string) Builder {
	return b
}

func (b *brokerClientConfig) Build() BrokerClient {
	return nil
}

func NewBuilder() Builder {
	return &brokerClientConfig{}
}

type BrokerSocket rsocket.RSocket

type BrokerClient interface {
	AddService(rs rrpc.RequestHandlingRSocket) error

	AddNamedSocket(name string, rs rsocket.RSocket) error

	GroupServiceSocket(group string, tags Tags) BrokerSocket

	BroadcastServiceSocket(group string, tags Tags) BrokerSocket

	ShardServiceSocket(group string, shardKey []byte, tags Tags) BrokerSocket

	GroupNamedRSocket(group string, tags Tags) BrokerSocket

	BroadcastNamedRSocket(group string, tags Tags) BrokerSocket

	ShardNamedRSocket(group string, shardKey []byte, tags Tags) BrokerSocket
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
		return rsocket2.NewReconnectingRSocket(func() string {
			nodes := <-s.config.ds.DiscoverNodes()
			l := len(nodes)
			i := rand.Intn(l)
			return nodes[i]
		}, s.config.rh, func() payload.SetupPayload {
			return nil
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
	wrapper := rsocket2.NewNamedRSocketClientWrapper(name, rs)
	return b.config.rh.Register(wrapper)
}

func (b *brokerClient) GroupServiceSocket(group string, tags Tags) BrokerSocket {
	panic("implement me")
}

func (b *brokerClient) BroadcastServiceSocket(group string, tags Tags) BrokerSocket {
	panic("implement me")
}

func (b *brokerClient) ShardServiceSocket(group string, shardKey []byte, tags Tags) BrokerSocket {
	panic("implement me")
}

func (b *brokerClient) GroupNamedRSocket(group string, tags Tags) BrokerSocket {
	panic("implement me")
}

func (b *brokerClient) BroadcastNamedRSocket(group string, tags Tags) BrokerSocket {
	panic("implement me")
}

func (b *brokerClient) ShardNamedRSocket(group string, shardKey []byte, tags Tags) BrokerSocket {
	panic("implement me")
}

func (b *brokerClient) group(group string, tags Tags) BrokerSocket {
	socket := &rsocket2.PayloadTransformingRSocket{
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

func (b *brokerClient) broadcast(group string, tags Tags) BrokerSocket {
	socket := &rsocket2.PayloadTransformingRSocket{
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

func (b *brokerClient) shard(group string, shardKey []byte, tags Tags) BrokerSocket {
	socket := &rsocket2.PayloadTransformingRSocket{
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
