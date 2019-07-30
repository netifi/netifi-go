package framing

import (
	"bytes"
	"encoding/binary"
	"fmt"
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
		e = fmt.Errorf("error writing header: %s", e)
		return
	}

	if len(ipAddr) > 0 {
		e = binary.Write(w, binary.BigEndian, ipAddr)
		if e != nil {
			e = fmt.Errorf("error writing ip address: %s", e)
			return
		}
	} else {
		e = binary.Write(w, binary.BigEndian, uint32(0))
		if e != nil {
			return
		}
	}

	l := len(group)
	e = binary.Write(w, binary.BigEndian, uint32(l))
	if e != nil {
		e = fmt.Errorf("error writing group: %s", e)
		return
	}
	e = binary.Write(w, binary.BigEndian, []byte(group))
	if e != nil {
		e = fmt.Errorf("error writing group: %s", e)
		return
	}

	e = binary.Write(w, binary.BigEndian, key)
	if e != nil {
		e = fmt.Errorf("error writing access key: %s", e)
		return
	}

	e = binary.Write(w, binary.BigEndian, uint32(len(token)))
	if e != nil {
		e = fmt.Errorf("error writing access token: %s", e)
		return
	}

	e = binary.Write(w, binary.BigEndian, token)
	if e != nil {
		e = fmt.Errorf("error writing token: %s", e)
		return
	}

	e = binary.Write(w, binary.BigEndian, uuid)
	if e != nil {
		e = fmt.Errorf("error writing uuid: %s", e)
		return
	}

	e = binary.Write(w, binary.BigEndian, additionalFlags)
	if e != nil {
		e = fmt.Errorf("error writing additional flags: %s", e)
		return
	}

	for _, tag := range tags {
		key := tag.Key()
		value := tag.Value()

		e = binary.Write(w, binary.BigEndian, uint32(len(key)))
		if e != nil {
			e = fmt.Errorf("error writing tag key: %s", e)
			return
		}

		e = binary.Write(w, binary.BigEndian, []byte(key))
		if e != nil {
			e = fmt.Errorf("error writing tag key: %s", e)
			return
		}

		e = binary.Write(w, binary.BigEndian, uint32(len(value)))
		if e != nil {
			e = fmt.Errorf("error writing tag value: %s", e)
			return
		}

		e = binary.Write(w, binary.BigEndian, []byte(value))
		if e != nil {
			e = fmt.Errorf("error writing tag value: %s", e)
			return
		}
	}
	d = w.Bytes()
	return
}
