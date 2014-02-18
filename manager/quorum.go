package manager

import "github.com/campadrenalin/go-deje/model"

type QuorumManager struct {
	genericManager
}

func NewQuorumManager() *QuorumManager {
	om := newGenericManager()
	return &QuorumManager{om}
}

func (qm *QuorumManager) Register(quorum model.Quorum) {
	qm.register(quorum)
}

func (qm *QuorumManager) Unregister(quorum model.Quorum) {
	qm.unregister(quorum)
}

func (qm *QuorumManager) DeserializeFrom(items map[string]model.Quorum) {
	for _, value := range items {
		qm.Register(value)
	}
}

func (qm *QuorumManager) SerializeTo(items map[string]model.Quorum) {
	for key, value := range qm.GetItems() {
		ev := value.(model.Quorum)
		items[key] = ev
	}
}
