package manager

import "github.com/campadrenalin/go-deje/model"

type QuorumManager struct {
	GenericManager
}

func NewQuorumManager() QuorumManager {
	om := NewGenericManager()
	return QuorumManager{om}
}

func (qm *QuorumManager) Register(quorum model.Quorum) {
	qm.register(quorum)
}

func (qm *QuorumManager) Unregister(quorum model.Quorum) {
	qm.unregister(quorum)
}
