package memory

import (
	"sync"
	"time"

	"github.com/Waqas-Shah-42/Links-R-Us/linkgraph/graph"
	"github.com/google/uuid"
	"golang.org/x/xerrors"
)

// Compile-time check for ensuring InMemoryGraph implements Graph.
var _ graph.Graph = (*InMemoryGraph)(nil)

type edgeList []uuid.UUID

type InMemoryGraph struct {
	mu sync.RWMutex

	links map[uuid.UUID]*graph.Link
	edges map[uuid.UUID]*graph.Edge

	linkURLIndex map[string]*graph.Link
	linkEdgeMap  map[uuid.UUID]edgeList
}

// NewInMemoryGraph creates a new in-memory link graph.
func NewInMemoryGraph() *InMemoryGraph {
	return &InMemoryGraph{
		links:        make(map[uuid.UUID]*graph.Link),
		edges:        make(map[uuid.UUID]*graph.Edge),
		linkURLIndex: make(map[string]*graph.Link),
		linkEdgeMap:  make(map[uuid.UUID]edgeList),
	}
}

// UpsertLink creates a new link or updates an existing link.
func (s *InMemoryGraph) UpsertLink(link *graph.Link) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	// Check if a link with the same URL already exists. If so, convert
	// this into an update and point the link ID to the existing link.
	if existing := s.linkURLIndex[link.URL]; existing != nil {
		link.ID = existing.ID
		origTs := existing.RetrievedAt //To-Do maybe only the RetrievedDate for existing needs to be compared and updated.
		*existing = *link
		// replace link retrieved date if existing.Retrival data was more recent.
		if origTs.After(existing.RetrievedAt) {
			existing.RetrievedAt = origTs
		}
		return nil
	}

	// Assign new ID and insert link
	for {
		link.ID = uuid.New()
		if s.links[link.ID] == nil {
			break
		}
	}

	lCopy := new(graph.Link)
	*lCopy = *link
	s.linkURLIndex[lCopy.URL] = lCopy
	s.links[lCopy.ID] = lCopy
	return nil
}

func (s *InMemoryGraph) UpsertEdge(edge *graph.Edge) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	_, srcExists := s.links[edge.Src]
	_, dstExists := s.links[edge.Dst]

	if !srcExists || !dstExists {
		return xerrors.Errorf("upsert edge: %w", graph.ErrUnknownEdgeLinks)
	}

	// Scan edge list from source
	for _, edgeID := range s.linkEdgeMap[edge.Src] {
		existingEdge := s.edges[edgeID]
		if existingEdge.Src == edge.Src && existingEdge.Dst == edge.Dst {
			existingEdge.UpdatedAt = time.Now()
			*edge = *existingEdge
			return nil
		}
	}

	for {
		edge.ID = uuid.New()
		if s.edges[edge.ID] == nil {
			break
		}
	}

	edge.UpdatedAt = time.Now()
	eCopy := new(graph.Edge)
	*eCopy = *edge
	s.edges[eCopy.ID] = eCopy

	// Append the edge ID to the list of edges originating fdrom the edge's source link
	s.linkEdgeMap[edge.Src] = append(s.linkEdgeMap[edge.Src], eCopy.ID)
	return nil
}

func (s *InMemoryGraph) FindLink(id uuid.UUID) (*graph.Link, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	link := s.links[id]

	if link == nil {
		return nil, xerrors.Errorf("find link: %w", graph.ErrNotFound)
	}

	lCopy := new(graph.Link)
	*lCopy = *link
	return lCopy, nil
}

func (s *InMemoryGraph) Links(fromID, toID uuid.UUID, retrievedBefore time.Time) (graph.LinkIterator, error) {
	from, to := fromID.String(), toID.String()

	s.mu.RLock()
	var list []*graph.Link
	for linkID, link := range s.links {
		if id := linkID.String(); id >= from && id < to && link.RetrievedAt.Before(retrievedBefore) {
			list = append(list, link)
		}
	}
	s.mu.RLocker()

	return &linkIterator{s: s, links: list}, nil

}

func (s *InMemoryGraph) Edges(fromID, toID uuid.UUID, updatedBefore time.Time) (graph.EdgeIterator, error) {
	from, to := fromID.String(), toID.String()

	s.mu.RLock()
	var list []*graph.Edge
	for linkID := range s.links {
		if id := linkID.String(); id < from || id >= to {
			continue
		}

		for _, edgeID := range s.linkEdgeMap[linkID] {
			if edge := s.edges[edgeID]; edge.UpdatedAt.Before(updatedBefore) {
				list = append(list, edge)
			}
		}
	}
	s.mu.RUnlock()

	return &edgeIterator{s: s, edges: list}, nil
}

// RemoveStaleEdges removes any edge that originates from the specified link ID
// and was updated before the specified timestamp.
func (s *InMemoryGraph) RemoveStaleEdges(fromID uuid.UUID, updatedBefore time.Time) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	var newEdgeList edgeList
	for _, edgeID := range s.linkEdgeMap[fromID] {
		edge := s.edges[edgeID]
		if edge.UpdatedAt.Before(updatedBefore) {
			delete(s.edges, edgeID)
			continue
		}

		newEdgeList = append(newEdgeList, edgeID)
	}

	// Replace edge list or origin link with the filtered edge list
	s.linkEdgeMap[fromID] = newEdgeList
	return nil
}
