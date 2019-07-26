package discovery_strategy

type DiscoveryStrategy interface {
	DiscoverNodes() <-chan []string
}

type StaticDiscoveryStrategy struct {
	nodes []string
}

func (s *StaticDiscoveryStrategy) DiscoverNodes() <-chan []string {
	n := make(chan []string, 1)
	n <- s.nodes
	close(n)
	return n
}

func New(n ...string) *StaticDiscoveryStrategy {
	return &StaticDiscoveryStrategy{
		nodes: n,
	}
}
