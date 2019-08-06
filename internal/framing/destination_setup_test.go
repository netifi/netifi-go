package framing

import (
	"encoding/base64"
	"github.com/google/uuid"
	"github.com/netifi/netifi-go/tags"
	"github.com/stretchr/testify/assert"
	"net"
	"testing"
)

func TestEncodeWithoutInetAddressAndNoTags(t *testing.T) {
	bytes, err := base64.StdEncoding.DecodeString("kTBDVtfRBO4tHOnZzSyY5ym2kfY=")
	if err != nil {
		t.Error(err)
	}
	u, _ := uuid.NewRandom()
	d, err := EncodeDestinationSetup(nil, "group", 123, bytes, u, 0, tags.Empty())

	setup := DestinationSetup(d)
	assert.Equal(t, "group", setup.Group())
	assert.Equal(t, uint64(123), setup.AccessKey())
	assert.Equal(t, bytes, setup.AccessToken())
	assert.Equal(t, u, setup.ConnectionId())
	assert.Equal(t, uint16(0), setup.AdditionalFlags())
	assert.Equal(t, tags.Empty(), setup.Tags())
}

func TestEncodeWithoutInetAddressAndTags(t *testing.T) {
	bytes, err := base64.StdEncoding.DecodeString("kTBDVtfRBO4tHOnZzSyY5ym2kfY=")
	if err != nil {
		t.Error(err)
	}
	u, _ := uuid.NewRandom()
	of, _ := tags.Of("k1", "v1", "k2", "v2")
	d, err := EncodeDestinationSetup(nil, "group", 123, bytes, u, 0, of)

	setup := DestinationSetup(d)
	assert.Equal(t, "group", setup.Group())
	assert.Equal(t, uint64(123), setup.AccessKey())
	assert.Equal(t, bytes, setup.AccessToken())
	assert.Equal(t, u, setup.ConnectionId())
	assert.Equal(t, uint16(0), setup.AdditionalFlags())
	assert.Equal(t, of, setup.Tags())
}

func TestWithInetAddressAndNoTags(t *testing.T) {
	bytes, err := base64.StdEncoding.DecodeString("kTBDVtfRBO4tHOnZzSyY5ym2kfY=")
	if err != nil {
		t.Error(err)
	}
	u, _ := uuid.NewRandom()
	ip := net.IPv4(127, 0, 0, 1)
	d, err := EncodeDestinationSetup(ip, "group", 123, bytes, u, 0, tags.Empty())

	setup := DestinationSetup(d)
	assert.Equal(t, "group", setup.Group())
	assert.Equal(t, uint64(123), setup.AccessKey())
	assert.Equal(t, bytes, setup.AccessToken())
	assert.Equal(t, u, setup.ConnectionId())
	assert.Equal(t, uint16(0), setup.AdditionalFlags())
	assert.Equal(t, tags.Empty(), setup.Tags())
	assert.Equal(t, ip, setup.InetAddress())
}

func TestWithInetAddressAndTags(t *testing.T) {
	bytes, err := base64.StdEncoding.DecodeString("kTBDVtfRBO4tHOnZzSyY5ym2kfY=")
	if err != nil {
		t.Error(err)
	}
	u, _ := uuid.NewRandom()
	ip := net.IPv4(127, 0, 0, 1)
	of, _ := tags.Of("k1", "v1", "k2", "v2")
	d, err := EncodeDestinationSetup(ip, "group", 123, bytes, u, 0, of)

	setup := DestinationSetup(d)
	assert.Equal(t, "group", setup.Group())
	assert.Equal(t, uint64(123), setup.AccessKey())
	assert.Equal(t, bytes, setup.AccessToken())
	assert.Equal(t, u, setup.ConnectionId())
	assert.Equal(t, uint16(0), setup.AdditionalFlags())
	assert.Equal(t, of, setup.Tags())
	assert.Equal(t, ip, setup.InetAddress())

}
