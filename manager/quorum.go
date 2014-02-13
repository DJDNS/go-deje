package manager

import "github.com/campadrenalin/go-deje/model"

type QuorumManager struct {
	ObjectManager
}

func NewQuorumManager() QuorumManager {
	om := NewObjectManager()
	return QuorumManager{om}
}

func (qm *QuorumManager) Register(quorum model.Quorum) {
	qm.register(quorum)
}

func (qm *QuorumManager) Unregister(quorum model.Quorum) {
	qm.unregister(quorum)
}
