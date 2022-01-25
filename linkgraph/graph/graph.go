package graph

import (
	"time"

	"github.com/google/uuid"
)

type Iterator interface {
	// If next item doesn't exist or an error occours, Next() returns false.
	Next() bool

	Error() error

	Close() error
}

// Link encapsulates all data for the URL required by the project
type Link struct {
	ID          uuid.UUID
	URL         string
	RetrievedAt time.Time
}

// Edge describes a graph edge that originates from src and terminates at Dst
type Edge struct {
	ID        uuid.UUID
	Src       uuid.UUID
	Dst       uuid.UUID
	UpdatedAt time.Time
}

type LinkIterator interface {
	Iterator

	Link() *Link
}

type EdgeIterator interface {
	Iterator

	Edge() *Edge
}

type Graph interface {
	UpserLink(link *Link) error
	FindLink(link *Link) (*Link, error)

	UpsertEdge(edge *Edge) error
	RemoveStaleEdges(fromID uuid.UUID, updatedBefore time.Time) error

	Links(fromID, toID uuid.UUID, retrievedBefore time.Time) (LinkIterator, error)
	Edges(fromID, toID uuid.UUID, updatedBefore time.Time) (EdgeIterator, error)
}
