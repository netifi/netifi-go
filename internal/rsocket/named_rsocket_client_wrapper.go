package rsocket

import (
	"github.com/rsocket/rsocket-go"
	"github.com/rsocket/rsocket-go/payload"
	rrpc "github.com/rsocket/rsocket-rpc-go"
)

type NamedRSocketClientWrapper struct {
	PayloadTransformingRSocket
	rrpc.RrpcRSocket
	name string
}

func (r *NamedRSocketClientWrapper) Name() string {
	return r.name
}

func (r *NamedRSocketClientWrapper) Wrap(msg payload.Payload) (p payload.Payload, err error) {
	d := msg.Data()
	m, _ := msg.Metadata()
	md, err := rrpc.EncodeMetadata(r.Name(), r.Name(), nil, m)

	if err != nil {
		return
	}

	p = payload.New(d, md)
	return
}

func NewNamedRSocketClientWrapper(name string, source rsocket.RSocket) rrpc.RrpcRSocket {
	wrapper := &NamedRSocketClientWrapper{}
	wrapper.name = name
	wrapper.Transformer = wrapper.Wrap
	wrapper.Source = func() rsocket.RSocket {
		return source
	}
	return wrapper
}
