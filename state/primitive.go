package state

import "errors"

type Primitive interface {
	Apply(*DocumentState) error
	Reverse(*DocumentState) (Primitive, error)
}

// Given a root and a path (call the pointed-to Container X),
// return X's parent, and the key from the parent to X.
//
// This is used throughout the primitives, so it makes sense
// to implement it as common code.
func getTraversal(c Container, path []interface{}) (Container, interface{}, error) {
	if len(path) == 0 {
		return nil, nil, errors.New("Empty path - must have >= 1 key")
	}
	traversal, last := path[:len(path)-1], path[len(path)-1]
	parent, err := Traverse(c, traversal)
	if err != nil {
		return nil, nil, err
	} else {
		return parent, last, nil
	}
}
