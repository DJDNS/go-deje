package manager

import "github.com/campadrenalin/go-deje/model"

// Base Manager from which others are derived
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

// Get a map[string]Manageable of all items in this Manager.
func (om *genericManager) GetItems() model.ManageableSet {
	return om.by_key
}

// Get the number of items in this Manager.
func (om *genericManager) Length() int {
	return len(om.by_key)
}

// Test whether this Manager contains a specific item.
func (om *genericManager) Contains(m model.Manageable) bool {
	return om.by_key.Contains(m)
}

// Get a specific item by its key.
func (om *genericManager) GetByKey(key string) (model.Manageable, bool) {
	m, ok := om.by_key[key]
	return m, ok
}

// Get a set of items by their group key.
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
