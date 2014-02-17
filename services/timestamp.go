package services

import "github.com/campadrenalin/go-deje/model"

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
	GetAfter(dochash string, after model.BlockHeight)

	MakeTimestamp(dochash string, qhash string)
}

type DummyTimestampService struct{}

func (tss DummyTimestampService) GetAfter(dochash string, after model.BlockHeight) {
	return
}

func (tss DummyTimestampService) MakeTimestamp(dochash string, qhash string) {
	return
}
