package named_rsocket_wrapper

import (
	"context"
	"fmt"
	"github.com/rsocket/rsocket-go/payload"
	"github.com/rsocket/rsocket-go/rx"
	"github.com/rsocket/rsocket-go/rx/flux"
	"github.com/rsocket/rsocket-go/rx/mono"
	rrpc "github.com/rsocket/rsocket-rpc-go"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestGetName(t *testing.T) {
	wrapper := NewClientWrapper("name", nil)
	assert.Equal(t, "name", wrapper.Name())
}

func TestWrap(t *testing.T) {
	wrapper := NewClientWrapper("test", &testRSocket{})
	response := wrapper.RequestResponse(payload.NewString("data", "metadata"))
	c, e := mono.ToChannel(response, context.Background())
	if e == nil {
		t.Error(e)
	}
	i := <-c
	data := payload.Payload(i).DataUTF8()
	assert.Equal(t, "data", data)
	metadata, ok := payload.Payload(i).Metadata()
	assert.True(t, ok)
	md := rrpc.Metadata(metadata)
	fmt.Println(md.String())
	assert.Equal(t, "test", md.Method())
	assert.Equal(t, "test", md.Service())
	innerFrame := string(md.Metadata())
	assert.Equal(t, "metadata", innerFrame)
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
	panic("implement me")
}

func (testRSocket) RequestChannel(msgs rx.Publisher) flux.Flux {
	panic("implement me")
}
