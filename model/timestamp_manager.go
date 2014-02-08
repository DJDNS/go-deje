package model

type TimestampManager struct {
	ObjectManager
}

func NewTimestampManager() TimestampManager {
	om := NewObjectManager()
	return TimestampManager{om}
}

func (em *TimestampManager) Register(timestamp Timestamp) {
	em.register(timestamp)
}

func (em *TimestampManager) Unregister(timestamp Timestamp) {
	em.unregister(timestamp)
}
