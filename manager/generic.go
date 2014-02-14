package manager

import "github.com/campadrenalin/go-deje/model"

type genericManager struct {
	by_key   model.ManageableSet
	by_group map[string]model.ManageableSet
}

func newGenericManager() genericManager {
	return genericManager{
		by_key:   make(model.ManageableSet),
		by_group: make(map[string]model.ManageableSet),
	}
}

func (om *genericManager) GetItems() model.ManageableSet {
	return om.by_key
}

func (om *genericManager) Length() int {
	return len(om.by_key)
}

func (om *genericManager) Contains(m model.Manageable) bool {
	return om.by_key.Contains(m)
}

func (om *genericManager) GetByKey(key string) (model.Manageable, bool) {
	m, ok := om.by_key[key]
	return m, ok
}

func (om *genericManager) GetGroup(key string) model.ManageableSet {
	_, ok := om.by_group[key]
	if !ok {
		om.by_group[key] = make(model.ManageableSet)
	}
	return om.by_group[key]
}

func (om *genericManager) register(m model.Manageable) {
	k := m.GetKey()
	gk := m.GetGroupKey()
	group := om.GetGroup(gk)

	om.by_key[k] = m
	group[k] = m
}

func (om *genericManager) unregister(m model.Manageable) {
	k := m.GetKey()
	gk := m.GetGroupKey()
	group := om.GetGroup(gk)

	delete(om.by_key, k)
	delete(group, k)
}
