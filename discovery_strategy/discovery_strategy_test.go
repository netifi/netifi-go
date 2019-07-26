package discovery_strategy

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestStaticDiscoveryStrategy(t *testing.T) {
	ds := New("tcp://localhost:8001", "tcp://localhost:8002")
	strings := <-ds.DiscoverNodes()

	s1 := strings[0]
	s2 := strings[1]

	assert.Equal(t, "tcp://localhost:8001", s1)
	assert.Equal(t, "tcp://localhost:8002", s2)
}
