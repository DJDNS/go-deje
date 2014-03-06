package state

type Primitive interface {
	Apply(*DocumentState) error
	Reverse(*DocumentState) (Primitive, error)
}

// Corresponds to SET builtin event handler.
type SetPrimitive struct {
	Path  []interface{}
	Value interface{}
}

func (p *SetPrimitive) Apply(ds *DocumentState) error {
	if len(p.Path) == 0 {
		container, err := MakeContainer(p.Value)
		if err != nil {
			return err
		}
		ds.Value = container
		return nil
	}

	traversal := p.Path[:len(p.Path)-1]
	last := p.Path[len(p.Path)-1]
	parent, err := Traverse(ds.Value, traversal)
	if err != nil {
		return err
	}
	return parent.Set(last, p.Value)
}
