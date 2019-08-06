package transforming_rsocket

import (
	"github.com/rsocket/rsocket-go"
	"github.com/rsocket/rsocket-go/payload"
	"github.com/rsocket/rsocket-go/rx"
	"github.com/rsocket/rsocket-go/rx/flux"
	"github.com/rsocket/rsocket-go/rx/mono"
)

type PayloadTransformingRSocket struct {
	Transformer func(msg payload.Payload) (payload.Payload, error)
	Source      func() rsocket.RSocket
}

func (r *PayloadTransformingRSocket) FireAndForget(msg payload.Payload) {
	p, e := r.Transformer(msg)
	if e != nil {
		panic(e)
	}

	r.Source().FireAndForget(p)
}

func (r *PayloadTransformingRSocket) MetadataPush(msg payload.Payload) {
	p, e := r.Transformer(msg)
	if e != nil {
		panic(e)
	}

	r.Source().MetadataPush(p)
}

func (r *PayloadTransformingRSocket) RequestResponse(msg payload.Payload) mono.Mono {
	p, e := r.Transformer(msg)
	if e != nil {
		return mono.Error(e)
	}

	return r.Source().RequestResponse(p)
}

func (r *PayloadTransformingRSocket) RequestStream(msg payload.Payload) flux.Flux {
	p, e := r.Transformer(msg)
	if e != nil {
		return flux.Error(e)
	}

	return r.Source().RequestStream(p)
}

func (r *PayloadTransformingRSocket) RequestChannel(msgs rx.Publisher) flux.Flux {
	return flux.Clone(msgs).Map(func(in payload.Payload) payload.Payload {
		p, e := r.Transformer(in)
		if e != nil {
			panic(e)
		}

		return p
	})
}
