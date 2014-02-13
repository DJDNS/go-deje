package manager

import "github.com/campadrenalin/go-deje/model"

type Manageable model.Manageable
type ManageableSet map[string]Manageable

func (ms ManageableSet) Contains(m Manageable) bool {
	stored, ok := ms[m.GetKey()]
	return ok && stored.Eq(m)
}

type GenericManager struct {
	by_key   ManageableSet
	by_group map[string]ManageableSet
}

func NewGenericManager() GenericManager {
	return GenericManager{
		by_key:   make(ManageableSet),
		by_group: make(map[string]ManageableSet),
	}
}

func (om *GenericManager) GetItems() ManageableSet {
	return om.by_key
}

func (om *GenericManager) Length() int {
	return len(om.by_key)
}

func (om *GenericManager) Contains(m Manageable) bool {
	return om.by_key.Contains(m)
}

func (om *GenericManager) GetByKey(key string) (Manageable, bool) {
	m, ok := om.by_key[key]
	return m, ok
}

func (om *GenericManager) GetGroup(key string) ManageableSet {
	_, ok := om.by_group[key]
	if !ok {
		om.by_group[key] = make(ManageableSet)
	}
	return om.by_group[key]
}

func (om *GenericManager) register(m Manageable) {
	k := m.GetKey()
	gk := m.GetGroupKey()
	group := om.GetGroup(gk)

	om.by_key[k] = m
	group[k] = m
}

func (om *GenericManager) unregister(m Manageable) {
	k := m.GetKey()
	gk := m.GetGroupKey()
	group := om.GetGroup(gk)

	delete(om.by_key, k)
	delete(group, k)
}
