package model

// Common interface for Events and Quorums.
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
