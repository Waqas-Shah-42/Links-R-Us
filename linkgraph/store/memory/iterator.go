package memory

import "github.com/Waqas-Shah-42/Links-R-Us/linkgraph/graph"

// linkIterator is a graph.LinkIterator implementation for the in-memory graph.
type linkIterator struct {
	s *InMemoryGraph

	links    []*graph.Link
	curIndex int
}
