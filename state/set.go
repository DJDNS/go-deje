package state

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

	parent, last, err := getTraversal(ds.Value, p.Path)
	if err != nil {
		return err
	}
	return parent.SetChild(last, p.Value)
}

func (p *SetPrimitive) Reverse(ds *DocumentState) (Primitive, error) {
	primitive := &SetPrimitive{
		Path:  []interface{}{},
		Value: ds.Export(),
	}
	return primitive, nil
}
