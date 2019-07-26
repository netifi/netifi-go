package netifi_go

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestStaticDiscoveryStrategy(t *testing.T) {
	strategy := NewStaticDiscoveryStrategy("tcp://localhost:8001", "tcp://localhost:8002")
	strings := <-strategy.DiscoverNodes()

	s1 := strings[0]
	s2 := strings[1]

	assert.Equal(t, "tcp://localhost:8001", s1)
	assert.Equal(t, "tcp://localhost:8002", s2)
}