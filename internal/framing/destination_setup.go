package framing

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"github.com/google/uuid"
	"github.com/netifi/netifi-go/tags"
	"github.com/pkg/errors"
	"net"
)

type DestinationSetup []byte

func (d DestinationSetup) InetAddress() net.IP {
	offset := HeaderOffset()

	inetAddressLength := int(binary.BigEndian.Uint32(d[offset : offset+4]))
	offset += 4

	if inetAddressLength > 0 {
		i := d[offset : offset+inetAddressLength]
		return net.IP(i)
	} else {
		return nil
	}

}

func (d DestinationSetup) Group() string {
	offset := HeaderOffset()

	inetAddressLength := int(binary.BigEndian.Uint32(d[offset : offset+4]))
	offset += 4 + inetAddressLength

	groupLength := int(binary.BigEndian.Uint32(d[offset : offset+4]))
	offset += 4

	return string(d[offset : offset+groupLength])
}

func (d DestinationSetup) AccessKey() uint64 {
	offset := HeaderOffset()

	inetAddressLength := int(binary.BigEndian.Uint32(d[offset : offset+4]))
	offset += 4 + inetAddressLength

	groupLength := int(binary.BigEndian.Uint32(d[offset : offset+4]))
	offset += 4 + groupLength

	return binary.BigEndian.Uint64(d[offset : offset+8])
}

func (d DestinationSetup) AccessToken() []byte {
	offset := HeaderOffset()

	inetAddressLength := int(binary.BigEndian.Uint32(d[offset : offset+4]))
	offset += 4 + inetAddressLength

	groupLength := int(binary.BigEndian.Uint32(d[offset : offset+4]))
	offset += 4 + groupLength + 8

	accessTokenLength := int(binary.BigEndian.Uint32(d[offset : offset+4]))
	offset += 4

	return d[offset : offset+accessTokenLength]
}

func (d DestinationSetup) ConnectionId() uuid.UUID {
	offset := HeaderOffset()

	inetAddressLength := int(binary.BigEndian.Uint32(d[offset : offset+4]))
	offset += 4 + inetAddressLength

	groupLength := int(binary.BigEndian.Uint32(d[offset : offset+4]))
	offset += 4 + groupLength + 8

	accessTokenLength := int(binary.BigEndian.Uint32(d[offset : offset+4]))
	offset += 4 + accessTokenLength

	u, e := uuid.FromBytes(d[offset : offset+16])
	if e != nil {
		panic(e)
	}
	return u
}

func (d DestinationSetup) AdditionalFlags() uint16 {
	offset := HeaderOffset()

	inetAddressLength := int(binary.BigEndian.Uint32(d[offset : offset+4]))
	offset += 4 + inetAddressLength

	groupLength := int(binary.BigEndian.Uint32(d[offset : offset+4]))
	offset += 4 + groupLength + 8

	accessTokenLength := int(binary.BigEndian.Uint32(d[offset : offset+4]))
	offset += 4 + accessTokenLength + 16

	return binary.BigEndian.Uint16(d[offset : offset+2])
}

func (d DestinationSetup) Tags() (t tags.Tags) {
	offset := HeaderOffset()

	inetAddressLength := int(binary.BigEndian.Uint32(d[offset : offset+4]))
	offset += 4 + inetAddressLength

	groupLength := int(binary.BigEndian.Uint32(d[offset : offset+4]))
	offset += 4 + groupLength + 8

	accessTokenLength := int(binary.BigEndian.Uint32(d[offset : offset+4]))
	offset += 4 + accessTokenLength + 16 + 2

	t = tags.Empty()
	for offset < len(d) {
		keyLength := int(binary.BigEndian.Uint32(d[offset : offset+4]))
		offset += 4

		key := string(d[offset : offset+keyLength])
		offset += keyLength

		valueLength := int(binary.BigEndian.Uint32(d[offset : offset+4]))
		offset += 4

		value := string(d[offset : offset+valueLength])
		offset += valueLength

		t = t.And(tags.New(key, value))
	}
	return
}

func EncodeDestinationSetup(ipAddr net.IP,
	group string, key uint64, token []byte,
	uuid uuid.UUID, additionalFlags uint16, tags tags.Tags) (d DestinationSetup, e error) {

	if token == nil {
		e = errors.New("access token is required")
		return
	}

	if len(group) == 0 {
		e = errors.New("group is required")
		return
	}

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
		e = binary.Write(w, binary.BigEndian, uint32(len(ipAddr)))
		if e != nil {
			e = fmt.Errorf("error writing ip address: %s", e)
			return
		}

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
	e = binary.Write(w, binary.BigEndian, int32(l))
	if e != nil {
		e = fmt.Errorf("error writing group: %s", e)
		return
	}
	_, e = w.WriteString(group)
	if e != nil {
		e = fmt.Errorf("error writing group: %s", e)
		return
	}

	e = binary.Write(w, binary.BigEndian, key)
	if e != nil {
		e = fmt.Errorf("error writing access key: %s", e)
		return
	}

	e = binary.Write(w, binary.BigEndian, int32(len(token)))
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

		e = binary.Write(w, binary.BigEndian, int32(len(key)))
		if e != nil {
			e = fmt.Errorf("error writing tag key: %s", e)
			return
		}

		_, e = w.WriteString(key)
		if e != nil {
			e = fmt.Errorf("error writing tag key: %s", e)
			return
		}

		e = binary.Write(w, binary.BigEndian, int32(len(value)))
		if e != nil {
			e = fmt.Errorf("error writing tag value: %s", e)
			return
		}

		_, e = w.WriteString(value)
		if e != nil {
			e = fmt.Errorf("error writing tag value: %s", e)
			return
		}
	}

	d = w.Bytes()
	return
}
