package model

// Represents a JSON map (object).
type JSONObject map[string]interface{}

// Interface for the go-deje.manager structs to store.
//
// We have to define it here, so that the model structs
// can implement Eq(Manageable).
type Manageable interface {
	GetKey() string
	GetGroupKey() string

	Eq(Manageable) bool
}
type ManageableSet map[string]Manageable

func (ms ManageableSet) Contains(m Manageable) bool {
	stored, ok := ms[m.GetKey()]
	return ok && stored.Eq(m)
}