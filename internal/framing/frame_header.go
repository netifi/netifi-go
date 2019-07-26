package framing

import (
	"bytes"
	"encoding/binary"
)

const (
	majorVersion = 0
	minorVersion = 1
	headerLength = 6
)

type Version [2]uint16
type FrameHeader []byte

func NewVersion(majorVersion uint16, minorVersion uint16) Version {
	return Version{majorVersion, minorVersion}
}

func (v Version) MajorVersion() uint16 {
	return v[0]
}

func (v Version) MinorVersion() uint16 {
	return v[1]
}

func (f FrameHeader) Version() Version {
	return Version{binary.BigEndian.Uint16(f[0:2]), binary.BigEndian.Uint16(f[2:4])}
}

func (f FrameHeader) FrameType() FrameType {
	return FrameType(binary.BigEndian.Uint16(f[4:6]))
}

func HeaderOffset() int {
	return headerLength
}

func encodeFrameHeader(frameType FrameType) (f FrameHeader, err error) {
	return encodeFrameHeaderWithVersion(NewVersion(majorVersion, minorVersion), frameType)
}

func encodeFrameHeaderWithVersion(version Version, frameType FrameType) (f FrameHeader, err error) {
	w := &bytes.Buffer{}
	err = binary.Write(w, binary.BigEndian, version[0])
	if err != nil {
		return
	}
	err = binary.Write(w, binary.BigEndian, version[1])
	if err != nil {
		return
	}
	err = binary.Write(w, binary.BigEndian, frameType)
	if err != nil {
		return
	}
	f = w.Bytes()
	return
}
