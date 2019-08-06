package netifi

import (
	"context"
	"fmt"
	"github.com/netifi/netifi-go/internal/rsocket/named_rsocket_wrapper"
	"github.com/netifi/netifi-go/internal/rsocket/unwrapping_rsocket"
	"github.com/netifi/netifi-go/tags"
	"github.com/rsocket/rsocket-go"
	"github.com/rsocket/rsocket-go/payload"
	"github.com/rsocket/rsocket-go/rx"
	"github.com/rsocket/rsocket-go/rx/flux"
	"github.com/rsocket/rsocket-go/rx/mono"
	rrpc "github.com/rsocket/rsocket-rpc-go"
	"github.com/stretchr/testify/assert"
	"strings"
	"testing"
)

func TestShouldIncludeUri(t *testing.T) {
	_, e := New().Build()
	if e == nil {
		t.Fail()
		return
	}
	fmt.Println(e)
	assert.True(t, strings.Contains(e.Error(), "uri or discovery service"))
}

func TestShouldIncludeAccessKeyGreaterThanZero(t *testing.T) {
	_, e := New().Uri("tcp://localhost").Build()
	if e == nil {
		t.Fail()
		return
	}
	fmt.Println(e)
	assert.True(t, strings.Contains(e.Error(), "key must be greater"))
}

func TestShouldIncludeAccessToken(t *testing.T) {
	_, e := New().AccessKey(123).Uri("tcp://localhost").Build()
	if e == nil {
		t.Fail()
		return
	}
	fmt.Println(e)
	assert.True(t, strings.Contains(e.Error(), "must include access token"))
}

func TestShouldIncludeGroup(t *testing.T) {
	_, e := New().
		AccessKey(123).
		AccessToken([]byte("kTBDVtfRBO4tHOnZzSyY5ym2kfY=")).
		Uri("tcp://localhost").
		Build()
	if e == nil {
		t.Fail()
		return
	}
	fmt.Println(e)
	assert.True(t, strings.Contains(e.Error(), "must include a group"))
}

func TestBuild(t *testing.T) {
	_, e := New().
		AccessKey(123).
		AccessToken([]byte("kTBDVtfRBO4tHOnZzSyY5ym2kfY=")).
		Uri("tcp://localhost:8001").
		Group("test").
		Build()
	if e != nil {
		fmt.Println(e)
		t.Fail()
		return
	}
}

func TestNamedSocket(t *testing.T) {
	client, e := New().Build()

	if e != nil {
		fmt.Println(e)
		t.Fail()
		return
	}

	gs := client.GroupNamedRSocket("mySocket", "test", tags.Empty())
	response := gs.RequestResponse(payload.NewString("some data", "some metadata"))
	payloads, errors := mono.ToChannel(response, context.Background())

	select {
	case p, ok := <-payloads:
		if ok {
			fmt.Println(p)
			data := payload.Payload(p).DataUTF8()
			metadata, b := payload.Payload(p).MetadataUTF8()
			if b {
				assert.Equal(t, "some metadata", metadata)
			} else {
				t.Fail()
			}

			assert.Equal(t, "some data", data)
		}
	case e := <-errors:
		if e != nil {
			fmt.Println(e)
			t.Fail()
		}
	}
}

type testRSocket struct {
}

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

type testRSocketSelector struct {
}

func (t testRSocketSelector) selectRSocket() rsocket.RSocket {
	handler := rrpc.NewRequestHandler()
	_ = handler.Register(named_rsocket_wrapper.NewServiceWrapper("mySocket", &testRSocket{}))
	return unwrapping_rsocket.New(handler)
}
