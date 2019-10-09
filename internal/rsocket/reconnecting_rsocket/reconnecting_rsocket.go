package reconnecting_rsocket

import (
	"context"
	"github.com/jjeffcaii/reactor-go/scheduler"
	"github.com/netifi/netifi-go/internal/rsocket/unwrapping_rsocket"
	"github.com/rsocket/rsocket-go"
	"github.com/rsocket/rsocket-go/payload"
	"github.com/rsocket/rsocket-go/rx"
	"github.com/rsocket/rsocket-go/rx/flux"
	"github.com/rsocket/rsocket-go/rx/mono"
	rrpc "github.com/rsocket/rsocket-rpc-go"
	"log"
	"sync"
	"time"
)

type ReconnectingRSocket struct {
	sync.Mutex
	requestHandler rrpc.RequestHandlingRSocket
	activeSocket   *rsocket.CloseableRSocket
	uri            func() string
	factory        func() payload.Payload
}

func New(uri func() string, requestHandler rrpc.RequestHandlingRSocket, factory func() payload.Payload) rsocket.RSocket {
	rs := &ReconnectingRSocket{
		uri:            uri,
		requestHandler: requestHandler,
		activeSocket:   nil,
		factory:        factory,
	}

	rs.GetRSocket()

	return rs
}

func (r *ReconnectingRSocket) ConnectRSocket() (rs rsocket.CloseableRSocket, e error) {
	r.Lock()
	defer r.Unlock()
	if r.activeSocket == nil {
		log.Println("no active rsocket attempting to connect to ", r.uri())
		rs, e = rsocket.
			Connect().
			MetadataMimeType("application/octet-stream").
			DataMimeType("application/octet-stream").
			SetupPayload(r.factory()).
			Acceptor(
				func(socket rsocket.RSocket) rsocket.RSocket {
					log.Println("successfully connected to ", r.uri())
					return unwrapping_rsocket.New(r.requestHandler)
				}).
			Transport(r.uri()).
			Start(context.Background())

		if e != nil {
			log.Printf("error connecting to uri %s{} - %s\n", r.uri(), e)
			return
		}

		rs.OnClose(r.Reset)

		r.activeSocket = &rs
	} else {
		rs = *r.activeSocket
	}

	return
}

func (r *ReconnectingRSocket) GetViaChannel(ctx context.Context) (rs <-chan rsocket.CloseableRSocket) {
	out := make(chan rsocket.CloseableRSocket)
	scheduler.Parallel().Worker().Do(func() {
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

func (r *ReconnectingRSocket) Reset(err error) {
	r.Lock()
	defer r.Unlock()
	log.Printf("connection reset to uri %s due to error %s", r.uri(), err.Error())
	r.activeSocket = nil
}

func (r *ReconnectingRSocket) FireAndForget(msg payload.Payload) {
	r.GetRSocket().FireAndForget(msg)
}

func (r *ReconnectingRSocket) MetadataPush(msg payload.Payload) {
	r.GetRSocket().MetadataPush(msg)
}

func (r *ReconnectingRSocket) RequestResponse(msg payload.Payload) mono.Mono {
	return r.GetRSocket().RequestResponse(msg)
}

func (r *ReconnectingRSocket) RequestStream(msg payload.Payload) flux.Flux {
	return r.GetRSocket().RequestStream(msg)
}

func (r *ReconnectingRSocket) RequestChannel(msgs rx.Publisher) flux.Flux {
	return r.GetRSocket().RequestChannel(msgs)
}
