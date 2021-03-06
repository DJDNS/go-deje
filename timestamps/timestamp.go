// Retrieval and processing of timestamps.
//
// Includes a few simple implementations of the TimestampService
// interface, so that not every test needs the full setup work of,
// say, the BitcoinTimestampService (not yet implemented).
package timestamps

import (
	"sort"

	"github.com/DJDNS/go-deje/document"
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
	GetTimestamps() ([]string, error)
}

// Always successfully returns an empty timestamp list.
type DummyTimestampService struct{}

func (tss DummyTimestampService) GetTimestamps() ([]string, error) {
	return make([]string, 0), nil
}

// A timestamp service that includes a Document pointer, and always
// returns the sorted list of event hashes when you call GetTimestamps.
//
// This is a useful approximation of real behavior for network-free
// testing, because the number of events will be small, and real
// timestamp services use hash sorting in any situation where the exact
// timing between two timestamps is ambiguous (multiple timestamps in the
// same block).
type SortingTimestampService struct {
	Doc document.Document
}

func NewSortingTimestampService(doc document.Document) SortingTimestampService {
	return SortingTimestampService{doc}
}
func (sts SortingTimestampService) GetTimestamps() ([]string, error) {
	items := sts.Doc.Events
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

// A timestamp service that includes a Document pointer, and always
// returns the doc.Timestamps list.
//
// Bandaid implementation for timestamps accumulated from peers.
type PeerTimestampService struct {
	Doc *document.Document
}

func NewPeerTimestampService(doc *document.Document) PeerTimestampService {
	return PeerTimestampService{doc}
}
func (sts PeerTimestampService) GetTimestamps() ([]string, error) {
	return sts.Doc.Timestamps, nil
}
