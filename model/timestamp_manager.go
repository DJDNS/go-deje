package model

type TimestampManager struct {
	ObjectManager
}

func NewTimestampManager() TimestampManager {
	om := NewObjectManager()
	return TimestampManager{om}
}

func (tm *TimestampManager) Register(timestamp Timestamp) {
	tm.register(timestamp)
}

func (tm *TimestampManager) Unregister(timestamp Timestamp) {
	tm.unregister(timestamp)
}
