package framing

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestVersion(t *testing.T) {
	version := NewVersion(0, 1)

	assert.Equal(t, uint16(0), version.MajorVersion())
	assert.Equal(t, uint16(1), version.MinorVersion())
}

func TestFrameHeader(t *testing.T) {
	version := NewVersion(0, 1)

	f, err := encodeFrameHeaderWithVersion(version, FrameTypeDestinationSetup)
	if err != nil {
		t.Fail()
	}

	v := f.Version()

	assert.Equal(t, uint16(0), v.MajorVersion())
	assert.Equal(t, uint16(1), v.MinorVersion())

	assert.Equal(t, FrameTypeDestinationSetup, f.FrameType())

}