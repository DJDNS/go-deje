package manager

import "github.com/campadrenalin/go-deje/model"

type QuorumManager struct {
	genericManager
}

func NewQuorumManager() QuorumManager {
	om := newGenericManager()
	return QuorumManager{om}
}

func (qm *QuorumManager) Register(quorum model.Quorum) {
	qm.register(quorum)
}

func (qm *QuorumManager) Unregister(quorum model.Quorum) {
	qm.unregister(quorum)
}
