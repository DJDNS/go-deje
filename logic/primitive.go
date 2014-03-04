package logic

// Represents a basic document state transformation.
// All Events are conceptually a list of primitives.
type EventPrimitive interface {
	Apply(*DocumentState) error
	GetReversal(*DocumentState) error
}

// Set the value of an object at the given path.
type SetPrimitive struct {
	Path  []interface{}
	Value interface{}
}

func (p *SetPrimitive) Apply(ds *DocumentState) error {
	return nil
}

// Delete the object at the given path.
type DeletePrimitive struct {
	Path []string
}
