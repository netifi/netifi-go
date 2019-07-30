package named_rsocket_wrapper

import (
	"github.com/netifi/netifi-go/internal/rsocket/transforming_rsocket"
	"github.com/rsocket/rsocket-go"
	"github.com/rsocket/rsocket-go/payload"
	rrpc "github.com/rsocket/rsocket-rpc-go"
)

type NamedRSocketServiceWrapper struct {
	transforming_rsocket.PayloadTransformingRSocket
	name string
}

func (r *NamedRSocketServiceWrapper) Name() string {
	return r.name
}

func (r *NamedRSocketServiceWrapper) wrap(msg payload.Payload) (p payload.Payload, err error) {
	d := msg.Data()
	m, _ := msg.Metadata()
	metadata := rrpc.Metadata(m).Metadata()
	p = payload.New(d, metadata)
	return
}

func NewServiceWrapper(name string, source rsocket.RSocket) rrpc.RrpcRSocket {
	wrapper := &NamedRSocketServiceWrapper{}
	wrapper.name = name
	wrapper.Transformer = wrapper.wrap
	wrapper.Source = func() rsocket.RSocket {
		return source
	}
	return wrapper
}