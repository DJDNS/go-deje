package deje

type Manageable interface {
	GetKey() string
	GetGroupKey() string
}

type ManageableSet map[string]Manageable

type ObjectManager struct {
	by_key   ManageableSet
	by_group map[string]ManageableSet
}

func (om *ObjectManager) Register(m Manageable) {
	k := m.GetKey()
	gk := m.GetGroupKey()
	group := om.by_group[gk]

	om.by_key[k] = m
	group[k] = m
}

func (om *ObjectManager) Unregister(m Manageable) {
	k := m.GetKey()
	gk := m.GetGroupKey()
	group := om.by_group[gk]

	delete(om.by_key, k)
	delete(group, k)
}
