package deje

type Timestamp struct {
	SyncHash    string
	BlockHeight uint64
}

func (t Timestamp) WasBefore(other Timestamp) bool {
	if t.BlockHeight == other.BlockHeight {
		return t.SyncHash < other.SyncHash
	} else {
		return t.BlockHeight < other.BlockHeight
	}
}
