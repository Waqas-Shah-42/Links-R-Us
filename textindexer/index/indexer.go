package indexer

import (
	"time"

	"github.com/google/uuid"
)


type Indexer interface {
	Index(doc *Document) error
	FindByID (linkID uuid.UUID) (*Document, error)
	Search (query Query) (Iterator, error)
 }

type Document struct {
	LinkID uuid.UUID

	URL string

	Title string
	Content string

	IndexedAt time.Time
	PageRank float64
}

type Query struct {
	Type QueryType
	Expression string
	Offset uint64
}

type QueryType uint8

const (
	QueryTypeMatch QueryType = iota
	QueryTypePhase
)

type Iterator interface {
	// close iterator
	Close() error

	// Loads the next documnet matching search query.
	// Return false if no more documnets left
	Next() bool

	// Returns last error encountered by iterator
	Error() error

	// Returns current document
	Document() *Document

	// Returns approximate number of search results
	TotalCount() uint64
}

/*

How the iterator can be used
// 'docIt' is a search iterator
for docIt.Next() {
doc := docIt.Document()
// Do something with doc...
}

if err := docIt.Error(); err != nil {
// Handle error...
}

*/