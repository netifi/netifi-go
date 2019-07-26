package transforming_rsocket

import (
	"context"
	"github.com/rsocket/rsocket-go"
	"github.com/rsocket/rsocket-go/payload"
	"github.com/rsocket/rsocket-go/rx"
)

type PayloadTransformingRSocket struct {
	rsocket.RSocket
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

func (r *PayloadTransformingRSocket) RequestResponse(msg payload.Payload) rx.Mono {
	p, e := r.Transformer(msg)
	if e != nil {
		return rx.NewMono(func(ctx context.Context, sink rx.MonoProducer) {
			sink.Error(e)
		})
	}

	return r.Source().RequestResponse(p)
}

func (r *PayloadTransformingRSocket) RequestStream(msg payload.Payload) rx.Flux {
	p, e := r.Transformer(msg)
	if e != nil {
		return rx.NewFluxFromError(e)
	}

	return r.Source().RequestStream(p)
}

func (r *PayloadTransformingRSocket) RequestChannel(msgs rx.Publisher) rx.Flux {
	background := context.Background()
	p1, e1 := rx.ToFlux(msgs).ToChannel(background)

	p2 := make(chan *payload.Payload)
	e2 := make(chan error)
	rx.ElasticScheduler().Do(background, func(ctx context.Context) {
		defer func() {
			close(p2)
			close(e2)
		}()
	loop:
		for {
			select {
			case s, ok := <-p1:
				if ok {
					p, e := r.Transformer(*s)
					if e != nil {
						p2 <- &p
					} else {
						e2 <- e
						break loop
					}
				} else {
					break loop
				}
			case e := <-e1:
				e2 <- e
				break loop
			}
		}
	})

	return r.Source().RequestChannel(rx.NewFluxFromChannel(p2, e2))
}
