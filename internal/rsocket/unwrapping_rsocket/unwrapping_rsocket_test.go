package unwrapping_rsocket

import (
	"context"
	"fmt"
	"github.com/netifi/netifi-go/internal/framing"
	"github.com/netifi/netifi-go/tags"
	"github.com/rsocket/rsocket-go/payload"
	"github.com/rsocket/rsocket-go/rx"
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
	payloads, errors := unwrapper.RequestResponse(p).ToChannel(context.Background())
	select {
	case p, ok := <-payloads:
		if ok {
			data := payload.Payload(*p).DataUTF8()
			assert.Equal(t, "data", data)
			metadata, b := payload.Payload(*p).MetadataUTF8()
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
	payloads, errors := unwrapper.RequestStream(p).ToChannel(context.Background())
	loop:
	for {
		select {
		case p, ok := <-payloads:
			if ok {
				data := payload.Payload(*p).DataUTF8()
				assert.Equal(t, "data", data)
				metadata, _ := payload.Payload(*p).MetadataUTF8()
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

func (testRSocket) RequestResponse(msg payload.Payload) rx.Mono {
	return rx.JustMono(msg)
}

func (testRSocket) RequestStream(msg payload.Payload) rx.Flux {
	data := msg.DataUTF8()
	metadata, _ := msg.MetadataUTF8()
	return rx.Range(0, 10).Map(func(n int) payload.Payload {

		fmt.Print("sending ")
		fmt.Println(n)
		newString := payload.NewString(data, fmt.Sprintf("%s - %d", metadata, n))
		fmt.Println(newString)
		return newString
	})
}

func (testRSocket) RequestChannel(msgs rx.Publisher) rx.Flux {
	panic("implement me")
}
