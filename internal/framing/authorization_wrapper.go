package framing

import (
	"bytes"
	"encoding/binary"
)

type AuthorizationWrapper []byte

func (a AuthorizationWrapper) AccessKey() uint64 {
	offset := HeaderOffset()
	return binary.BigEndian.Uint64(a[offset : offset+8])
}

func (a AuthorizationWrapper) InnerFrame() []byte {
	offset := HeaderOffset() + 8
	return a[offset:]
}

func encodeAuthorizationWrapper(accessKey uint64, innerFrame []byte) (a AuthorizationWrapper, err error) {
	w := &bytes.Buffer{}
	f, err := encodeFrameHeader(FrameTypeAuthorizationWrapper)
	if err != nil {
		return
	}
	err = binary.Write(w, binary.BigEndian, f)
	if err != nil {
		return
	}
	err = binary.Write(w, binary.BigEndian, accessKey)
	if err != nil {
		return
	}
	err = binary.Write(w, binary.BigEndian, innerFrame)
	if err != nil {
		return
	}
	a = w.Bytes()
	return
}
