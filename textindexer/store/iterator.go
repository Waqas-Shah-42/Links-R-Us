package memory

import (
	"github.com/Waqas-Shah-42/Links-R-Us/textindexer/index"
	"github.com/blevesearch/bleve"
)

type bleveIterator struct {
	idx       *InMemoryBleveIndexer
	searchReq *bleve.SearchRequest

	cumIdx uint64
	rsIdx  int
	rs     *bleve.SearchResult

	latchedDoc *index.Document
	lastErr    error
}

// Next loads next document matching query.
// Returns false if no more documents found
func (it *bleveIterator) Next() bool {
	if it.lastErr != nil || it.rs == nil || it.cumIdx >= it.rs.Total {
		return false
	}
	if it.rsIdx >= it.rs.Hits.Len() {
		it.searchReq.From += it.searchReq.Size
		if it.rs, it.lastErr = it.idx.idx.Search(it.searchReq); it.lastErr != nil {
			return false
		}
		it.rsIdx = 0
	}

	nextID := it.rs.Hits[it.rsIdx].ID
	if it.latchedDoc, it.lastErr = it.idx.findByID(nextID); it.lastErr != nil {
		return false
	}
	it.cumIdx++
	it.rsIdx++
	return true
}

// Close the iterator and release any allocated resources.
func (it *bleveIterator) Close() error {
	it.idx = nil
	it.searchReq = nil
	if it.rs != nil {
		it.cumIdx = it.rs.Total
	}
	return nil
}

// Error returns the last error encountered by the iterator.
func (it *bleveIterator) Error() error {
	return it.lastErr
}

// Document returns the current document from the result set.
func (it *bleveIterator) Document() *index.Document {
	return it.latchedDoc
}

// TotalCount returns the approximate number of search results.
func (it *bleveIterator) TotalCount() uint64 {
	if it.rs == nil {
		return 0
	}
	return it.rs.Total
}
