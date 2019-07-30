package unwrapping_rsocket

import (
	"errors"
	"fmt"
	"github.com/netifi/netifi-go/internal/framing"
	"github.com/netifi/netifi-go/internal/rsocket/transforming_rsocket"
	"github.com/rsocket/rsocket-go"
	"github.com/rsocket/rsocket-go/payload"
)

type UnwrappedRSocket struct {
	transforming_rsocket.PayloadTransformingRSocket
}

func (u *UnwrappedRSocket) unwrap(msg payload.Payload) (p payload.Payload, e error) {
	fmt.Print("before -> ")
	fmt.Println(msg)
	d := msg.Data()
	m, _ := msg.Metadata()
	header := framing.FrameHeader(m)
	frameType := header.FrameType()
	metadata, e := unwrapMetadata(frameType, m)
	p = payload.New(d, metadata)
	fmt.Print("after  -> ")
	fmt.Println(p)
	return
}

func unwrapMetadata(frameType framing.FrameType, m []byte) ([]byte, error) {
	switch frameType {
	case framing.FrameTypeAuthorizationWrapper:
		innerFrame := framing.AuthorizationWrapper(m).InnerFrame()
		return unwrapMetadata(framing.FrameHeader(innerFrame).FrameType(), innerFrame)
	case framing.FrameTypeBroadcast:
		return framing.Broadcast(m).Metadata(), nil
	case framing.FrameTypeGroup:
		return framing.Group(m).Metadata(), nil
	case framing.FrameTypeShard:
		return framing.Shard(m).Metadata(), nil
	default:
		return nil, errors.New("unknown metadata type")
	}
}

func New(source rsocket.RSocket) rsocket.RSocket {
	unwrapper := &UnwrappedRSocket{}
	unwrapper.Transformer = unwrapper.unwrap
	unwrapper.Source = func() rsocket.RSocket {
		return source
	}
	return unwrapper
}
