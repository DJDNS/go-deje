package model

import "github.com/campadrenalin/go-deje/serial"

type QuorumManager struct {
	ObjectManager
}

func NewQuorumManager() QuorumManager {
	om := NewObjectManager()
	return QuorumManager{om}
}

func (em *QuorumManager) Register(quorum Quorum) {
	em.register(quorum)
}

func (em *QuorumManager) Unregister(quorum Quorum) {
	em.unregister(quorum)
}

func (em *QuorumManager) DeserializeFrom(items map[string]serial.Quorum) {
	for _, value := range items {
		em.Register(QuorumFromSerial(value))
	}
}

func (em *QuorumManager) SerializeTo(items map[string]serial.Quorum) {
	for key, value := range em.GetItems() {
		ev := value.(Quorum)
		items[key] = ev.ToSerial()
	}
}
