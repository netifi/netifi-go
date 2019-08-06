package unwrapping_rsocket

import (
	"context"
	"fmt"
	"github.com/netifi/netifi-go/internal/framing"
	"github.com/netifi/netifi-go/tags"
	"github.com/rsocket/rsocket-go/payload"
	"github.com/rsocket/rsocket-go/rx"
	"github.com/rsocket/rsocket-go/rx/flux"
	"github.com/rsocket/rsocket-go/rx/mono"
	"github.com/stretchr/testify/assert"
	"strings"
	"testing"
)

func TestUnwrapping(t *testing.T) {
	unwrapper := New(&testRSocket{})
	tags, e := tags.Of("key", "value")
	if e != nil {
		t.Error(e)
		return
	}

	md, e := framing.EncodeGroup("testGroup", []byte("metadata"), tags)
	if e != nil {
		t.Error(e)
		return
	}

	p := payload.New([]byte("data"), md)
	response := unwrapper.RequestResponse(p)
	payloads, errors := mono.ToChannel(response, context.Background())
	select {
	case p, ok := <-payloads:
		if ok {
			data := payload.Payload(p).DataUTF8()
			assert.Equal(t, "data", data)
			metadata, b := payload.Payload(p).MetadataUTF8()
			if !b {
				t.Fail()
			}
			assert.Equal(t, "metadata", metadata)
		} else {
			t.Fail()
		}
	case e := <-errors:
		t.Error(e)
		t.Fail()
	}

}

func TestUnWrappingStream(t *testing.T) {
	unwrapper := New(&testRSocket{})
	tags, e := tags.Of("key", "value")
	if e != nil {
		t.Error(e)
		return
	}

	md, e := framing.EncodeGroup("testGroup", []byte("metadata"), tags)
	if e != nil {
		t.Error(e)
		return
	}

	p := payload.New([]byte("data"), md)
	payloads, errors := flux.ToChannel(unwrapper.RequestStream(p), context.Background())
loop:
	for {
		select {
		case p, ok := <-payloads:
			if ok {
				data := payload.Payload(p).DataUTF8()
				assert.Equal(t, "data", data)
				metadata, _ := payload.Payload(p).MetadataUTF8()
				fmt.Print("checking -> ")
				fmt.Println(metadata)
				split := strings.Split(metadata, " - ")
				assert.Equal(t, "metadata", split[0])
			} else {
				break loop
			}
		case e := <-errors:
			if e != nil {
				fmt.Println(e)
				t.Error(e)
				t.Fail()
				return
			}
		}
	}
}

type testRSocket struct{}

func (testRSocket) FireAndForget(msg payload.Payload) {
	panic("implement me")
}

func (testRSocket) MetadataPush(msg payload.Payload) {
	panic("implement me")
}

func (testRSocket) RequestResponse(msg payload.Payload) mono.Mono {
	return mono.Just(msg)
}

func (testRSocket) RequestStream(msg payload.Payload) flux.Flux {
	data := msg.DataUTF8()
	metadata, _ := msg.MetadataUTF8()

	return flux.Create(func(ctx context.Context, s flux.Sink) {
		go func() {
			for i := 0; i < 10; i++ {
				fmt.Print("sending ")
				fmt.Println(i)
				newString := payload.NewString(data, fmt.Sprintf("%s - %d", metadata, i))
				fmt.Println(newString)
			}
			s.Complete()
		}()
	})
}

func (testRSocket) RequestChannel(msgs rx.Publisher) flux.Flux {
	panic("implement me")
}
