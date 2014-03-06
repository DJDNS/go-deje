package state

type Primitive interface {
	Apply(*DocumentState) error
	Reverse(*DocumentState) (Primitive, error)
}

// Corresponds to SET builtin event handler.
type SetPrimitive
