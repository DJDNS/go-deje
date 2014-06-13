package timestamps

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

type DummyTimestampService struct{}

func (tss DummyTimestampService) GetTimestamps(topic string) ([]string, error) {
	return make([]string, 0), nil
}
