package reconnecting_rsocket

import (
	"context"
	"github.com/netifi/netifi-go/internal/rsocket/unwrapping_rsocket"
	"github.com/rsocket/rsocket-go"
	"github.com/rsocket/rsocket-go/payload"
	"github.com/rsocket/rsocket-go/rx"
	rrpc "github.com/rsocket/rsocket-rpc-go"
	"log"
	"sync"
	"time"
)

type ReconnectingRSocket struct {
	rsocket.RSocket
	mutex          sync.Mutex
	requestHandler rrpc.RequestHandlingRSocket
	activeSocket   rsocket.CloseableRSocket
	uri            func() string
	factory        func() payload.Payload
}

func New(uri func() string, requestHandler rrpc.RequestHandlingRSocket, factory func() payload.Payload) rsocket.RSocket {
	return &ReconnectingRSocket{
		uri:            uri,
		requestHandler: requestHandler,
		mutex:          sync.Mutex{},
		activeSocket:   nil,
		factory:        factory,
	}
}

func (r *ReconnectingRSocket) ConnectRSocket() (rs rsocket.CloseableRSocket, e error) {
	r.mutex.Lock()
	defer r.mutex.Unlock()
	if r.activeSocket == nil {
		rs, e = rsocket.
			Connect().
			MetadataMimeType("application/octet-stream").
			DataMimeType("application/octet-stream").
			SetupPayload(r.factory()).
			Acceptor(
				func(socket rsocket.RSocket) rsocket.RSocket {
					return unwrapping_rsocket.New(r.requestHandler)
				}).
			Transport(r.uri()).
			Start(context.Background())

		if e != nil {
			log.Printf("error connecting to uri %s{} - %s\n", r.uri(), e)
			return
		}

		rs.OnClose(r.Reset)

		r.activeSocket = rs
	}

	return
}

func (r *ReconnectingRSocket) GetViaChannel(ctx context.Context) (rs <-chan rsocket.CloseableRSocket) {
	out := make(chan rsocket.CloseableRSocket)
	rx.ElasticScheduler().Do(ctx, func(ctx context.Context) {
		defer close(out)

		for {
			rSocket, e := r.ConnectRSocket()
			if e == nil {
				out <- rSocket
				break
			} else {
				time.Sleep(1 * time.Second)
			}
		}

	})
	return out
}

func (r *ReconnectingRSocket) GetRSocket() rsocket.CloseableRSocket {
	rs, e := r.ConnectRSocket()
	if e != nil {
		out := r.GetViaChannel(context.Background())
		return <-out
	}
	return rs
}

func (r *ReconnectingRSocket) Reset() {
	r.mutex.Lock()
	defer r.mutex.Unlock()
	r.activeSocket = nil
}

func (r *ReconnectingRSocket) FireAndForget(msg payload.Payload) {
	r.GetRSocket().FireAndForget(msg)
}

func (r *ReconnectingRSocket) MetadataPush(msg payload.Payload) {
	r.GetRSocket().MetadataPush(msg)
}

func (r *ReconnectingRSocket) RequestResponse(msg payload.Payload) rx.Mono {
	return r.GetRSocket().RequestResponse(msg)
}

func (r *ReconnectingRSocket) RequestStream(msg payload.Payload) rx.Flux {
	return r.GetRSocket().RequestStream(msg)
}

func (r *ReconnectingRSocket) RequestChannel(msgs rx.Publisher) rx.Flux {
	return r.GetRSocket().RequestChannel(msgs)
}
