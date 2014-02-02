package deje

type BlockHeight uint64

type Timestamp struct {
	QHash string
	Time  BlockHeight
}

type TimestampSet map[string]*Timestamp

type TimestampManager struct {
	Stamps   TimestampSet
	PerBlock map[BlockHeight]TimestampSet
	LastPoll BlockHeight
}

func (t Timestamp) WasBefore(other Timestamp) bool {
	if t.Time == other.Time {
		return t.QHash < other.QHash
	} else {
		return t.Time < other.Time
	}
}

func (ts TimestampSet) Contains(t *Timestamp) bool {
	return ts[t.QHash] == t
}

func NewTimestampManager() TimestampManager {
	return TimestampManager{
		Stamps:   make(TimestampSet),
		PerBlock: make(map[BlockHeight]TimestampSet),
		LastPoll: 0,
	}
}

func (m *TimestampManager) GetBlock(time BlockHeight) TimestampSet {
	block, ok := m.PerBlock[time]
	if !ok {
		block = make(TimestampSet)
		m.PerBlock[time] = block
	}

	return block
}

func (m *TimestampManager) Register(ts *Timestamp) {
	hash := ts.QHash
	m.Stamps[hash] = ts

	block := m.GetBlock(ts.Time)
	block[hash] = ts
}

func (m *TimestampManager) Unregister(ts *Timestamp) {
	hash := ts.QHash
	delete(m.Stamps, hash)

	block := m.GetBlock(ts.Time)
	delete(block, hash)
}