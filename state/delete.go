package state

import "errors"

// Corresponds to DELETE builtin event handler.
type DeletePrimitive struct {
	Path []interface{}
}

func (p *DeletePrimitive) Apply(ds *DocumentState) error {
	if len(p.Path) == 0 {
		return errors.New("Cannot delete root node")
	}

	parent, last, err := getTraversal(ds.Value, p.Path)
	if err != nil {
		return err
	}
	return parent.RemoveChild(last)
}

func (p *DeletePrimitive) Reverse(ds *DocumentState) (Primitive, error) {
	primitive := &SetPrimitive{
		Path:  []interface{}{},
		Value: ds.Export(),
	}
	return primitive, nil
}
