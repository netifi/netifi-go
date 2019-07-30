package framing

import (
	"bytes"
	"encoding/binary"
	"github.com/google/uuid"
	"github.com/netifi/netifi-go/tags"
	"net"
)

type DestinationSetup []byte

func EncodeDestinationSetup(ipAddr net.IP,
	group string, key uint64, token []byte,
	uuid uuid.UUID, additionalFlags uint32, tags tags.Tags) (d DestinationSetup, e error) {
	w := &bytes.Buffer{}
	f, e := encodeFrameHeader(FrameTypeDestinationSetup)
	if e != nil {
		return
	}

	e = binary.Write(w, binary.BigEndian, f)
	if e != nil {
		return
	}

	if len(ipAddr) > 0 {
		e = binary.Write(w, binary.BigEndian, ipAddr)
		if e != nil {
			return
		}
	} else {
		e = binary.Write(w, binary.BigEndian, 0)
		if e != nil {
			return
		}
	}

	l := len(group)
	e = binary.Write(w, binary.BigEndian, uint32(l))
	if e != nil {
		return
	}
	e = binary.Write(w, binary.BigEndian, group)
	if e != nil {
		return
	}

	e = binary.Write(w, binary.BigEndian, key)
	if e != nil {
		return
	}

	e = binary.Write(w, binary.BigEndian, len(token))
	if e != nil {
		return
	}

	e = binary.Write(w, binary.BigEndian, token)
	if e != nil {
		return
	}

	e = binary.Write(w, binary.BigEndian, uuid)
	if e != nil {
		return
	}

	e = binary.Write(w, binary.BigEndian, additionalFlags)
	if e != nil {
		return
	}

	for _, tag := range tags {
		key := tag.Key()
		value := tag.Value()

		e = binary.Write(w, binary.BigEndian, uint32(len(key)))
		if e != nil {
			return
		}

		e = binary.Write(w, binary.BigEndian, key)
		if e != nil {
			return
		}

		e = binary.Write(w, binary.BigEndian, uint32(len(value)))
		if e != nil {
			return
		}

		e = binary.Write(w, binary.BigEndian, value)
		if e != nil {
			return
		}
	}
	d = w.Bytes()
	return
}
