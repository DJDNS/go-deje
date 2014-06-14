package state

import (
	"errors"
	"reflect"
)

// Represents a value in the tracked document state.
//
// Its most interesting attributes, from the point of view of external
// API, are that it can be Traversed and Exported.
type Container interface {
	GetChild(interface{}) (Container, error)
	SetChild(key, value interface{}) error
	RemoveChild(interface{}) error

	Export() interface{}
}

// Create a new container, based on the given object.
func makeContainer(value interface{}) (Container, error) {
	// Special case, since reflect.TypeOf(nil) == nil,
	// and nil.Kind() is a surefire recipe for runtime panic :/
	if reflect.TypeOf(value) == nil {
		return makeScalarContainer(value)
	}
	switch reflect.TypeOf(value).Kind() {
	case reflect.Map:
		as_map, ok := value.(map[string]interface{})
		if !ok {
			return nil, errors.New("Cannot cast map to map[string]interface{}")
		}
		return makeMapContainer(as_map)
	case reflect.Slice:
		as_slice, ok := value.([]interface{})
		if !ok {
			return nil, errors.New("Cannot cast slice to []interface{}")
		}
		return makeSliceContainer(as_slice)
	case reflect.Bool, reflect.Int, reflect.Uint, reflect.String:
		return makeScalarContainer(value)
	default:
		return nil, errors.New("Invalid type for containing")
	}
}

// Find a child based on a series of child keys.
// Will return an error for bad key types, unset keys, etc.
func Traverse(c Container, keys []interface{}) (Container, error) {
	var err error
	for _, key := range keys {
		c, err = c.GetChild(key)
		if err != nil {
			return nil, err
		}
	}
	return c, nil
}
