// Retrieval and processing of timestamps.
//
// Includes a few simple implementations of the TimestampService
// interface, so that not every test needs the full setup work of,
// say, the BitcoinTimestampService (not yet implemented).
package timestamps

import (
	"github.com/campadrenalin/go-deje/document"
	"sort"
)

// Different types of TimestampServices can be used,
// but the default implemetation makes use of the
// Bitcoin blockchain, given a bitcoind JSON-RPC location.
// Alternative timestamping services can simply be alternative
// blockchains, such as Litecoin, or something more exotic
// (if it can be made to fit the TimestampService interface).
//
// See https://en.bitcoin.it/wiki/API_reference_%28JSON-RPC%29
// for more information about this API.
type TimestampService interface {
	GetTimestamps(topic string) ([]string, error)
}

// Always successfully returns an empty timestamp list.
type DummyTimestampService struct{}

func (tss DummyTimestampService) GetTimestamps(topic string) ([]string, error) {
	return make([]string, 0), nil
}

// A timestamp service that includes a Document pointer, and always
// returns the sorted list of quorum hashes when you call GetTimestamps.
//
// This is a useful approximation of real behavior for network-free
// testing, because the number of quorums will be small, and real
// timestamp services use hash sorting in any situation where the exact
// timing between two timestamps is ambiguous (multiple timestamps in the
// same block).
type SortingTimestampService struct {
	Doc document.Document
}

func NewSortingTimestampService(doc document.Document) SortingTimestampService {
	return SortingTimestampService{doc}
}
func (sts SortingTimestampService) GetTimestamps(topic string) ([]string, error) {
	items := sts.Doc.Quorums
	timestamps := make([]string, len(items))

	// Extract keys as list
	var pos int
	for key := range items {
		timestamps[pos] = key
		pos++
	}

	// Sort and return
	sort.Strings(timestamps)
	return timestamps, nil
}
