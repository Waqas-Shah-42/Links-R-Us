package memory

import "github.com/blevesearch/bleve"

type bleveIterator struct {
	idx       *InMemoryBleveIndexer
	searchReq *bleve.SearchRequest
	cumIdx    uint64
	rsIdx     int
}
