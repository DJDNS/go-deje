package deje

import "strconv"

type BlockHeight uint64

type Timestamp struct {
	QHash string
	Time  BlockHeight
}

func (t Timestamp) GetKey() string {
	return t.QHash
}
func (t Timestamp) GetGroupKey() string {
	return strconv.FormatUint(uint64(t.Time), 10)
}
func (t Timestamp) Eq(other Manageable) bool {
	other_ts, ok := other.(Timestamp)
	if !ok {
		return false
	}
	return t == other_ts
}

func (t Timestamp) WasBefore(other Timestamp) bool {
	if t.Time == other.Time {
		return t.QHash < other.QHash
	} else {
		return t.Time < other.Time
	}
}
