package model

import "github.com/campadrenalin/go-deje/serial"

type QuorumManager struct {
	ObjectManager
}

func NewQuorumManager() QuorumManager {
	om := NewObjectManager()
	return QuorumManager{om}
}

func (qm *QuorumManager) Register(quorum Quorum) {
	qm.register(quorum)
}

func (qm *QuorumManager) Unregister(quorum Quorum) {
	qm.unregister(quorum)
}

func (qm *QuorumManager) DeserializeFrom(items map[string]serial.Quorum) {
	for _, value := range items {
		qm.Register(QuorumFromSerial(value))
	}
}

func (qm *QuorumManager) SerializeTo(items map[string]serial.Quorum) {
	for key, value := range qm.GetItems() {
		ev := value.(Quorum)
		items[key] = ev.ToSerial()
	}
}
