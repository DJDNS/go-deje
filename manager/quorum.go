package manager

import "github.com/campadrenalin/go-deje/model"

// Manager that stores Quorums.
//
// Quorums are grouped by their event hash, and keyed on their own hash.
type QuorumManager struct {
	genericManager
}

func NewQuorumManager() *QuorumManager {
	om := newGenericManager()
	return &QuorumManager{om}
}

// Register an Quorum. Afterwards, it can be retrieved by key or group.
func (qm *QuorumManager) Register(quorum model.Quorum) {
	qm.register(quorum)
}

// Remove an Quorum from the internal registries.
func (qm *QuorumManager) Unregister(quorum model.Quorum) {
	qm.unregister(quorum)
}

// Import a set of Quorums, registering all of them.
func (qm *QuorumManager) DeserializeFrom(items map[string]model.Quorum) {
	for _, value := range items {
		qm.Register(value)
	}
}

// Export a set of Quorums to a raw map.
func (qm *QuorumManager) SerializeTo(items map[string]model.Quorum) {
	for key, value := range qm.GetItems() {
		ev := value.(model.Quorum)
		items[key] = ev
	}
}
