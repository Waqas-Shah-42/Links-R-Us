package graph

import (
	"time"

	"github.com/google/uuid"
)

type Link struct {
	ID          uuid.UUID
	URL         string
	RetrievedAt time.Time
}

type Edge struct {
	ID        uuid.UUID
	Src       uuid.UUID
	Dst       uuid.UUID
	updatedAt time.Time
}

type Graph interface {
	UpserLink(link *Link) error
	FindLink(link *Link) (*Link, error)

	UpsertEdge(edge *Edge) error
	RemoveStaleEdges(fromID uuid.UUID, updatedBefore time.Time) error

	Links(fromID, toID uuid.UUID, retrievedBefore time.Time) (LinkIterator, error)
	Edges(fromID, toID uuid.UUID, updatedBefore time.Time) (EdgeIterator, error)
}
