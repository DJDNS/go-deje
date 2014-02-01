package deje

type Manageable interface {
	GetId() string
	GetParentId() string

	ProvenWrong() bool
	SetProvenWrong()
}

type ManageableSet map[string]Manageable

type ObjectManager struct {
	by_id     ManageableSet
	by_parent map[string]ManageableSet
}

func (om *ObjectManager) Register(m Manageable) {
	id := m.GetId()
	om.by_id[id] = m

	pid := m.GetParentId()
	pchildren := om.by_parent[pid]
	pchildren[id] = m
}

func (om *ObjectManager) Unregister(m Manageable) {
	id := m.GetId()
	pid := m.GetParentId()

	delete(om.by_id, id)
	pchildren := om.by_parent[pid]
	delete(pchildren, id)
}

func (om *ObjectManager) GetParent(m Manageable) (Manageable, bool) {
	pid := m.GetParentId()
	p, ok := om.by_id[pid]
	return p, ok
}

func (om *ObjectManager) GetRoot(m Manageable) Manageable {
	for {
		parent, had_p := om.GetParent(m)
		if had_p {
			m = parent
		} else {
			return m
		}
	}
}
