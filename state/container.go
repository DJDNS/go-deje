package state

import (
	"errors"
	"reflect"
)

type Container interface {
	GetChild(interface{}) (Container, error)
	Remove() error
	RemoveChild(interface{}) error
	SetParentage(Container, interface{})

	Set(key, value interface{}) error
	Export() interface{}
}

func MakeContainer(value interface{}) (Container, error) {
	// Special case, since reflect.TypeOf(nil) == nil,
	// and nil.Kind() is a surefire recipe for runtime panic :/
	if reflect.TypeOf(value) == nil {
		return MakeScalarContainer(value)
	}
	switch reflect.TypeOf(value).Kind() {
	case reflect.Map:
		as_map, ok := value.(map[string]interface{})
		if !ok {
			return nil, errors.New("Cannot cast map to map[string]interface{}")
		}
		return MakeMapContainer(as_map)
	case reflect.Slice:
		as_slice, ok := value.([]interface{})
		if !ok {
			return nil, errors.New("Cannot cast slice to []interface{}")
		}
		return MakeSliceContainer(as_slice)
	case reflect.Bool, reflect.Int, reflect.Uint, reflect.String:
		return MakeScalarContainer(value)
	default:
		return nil, errors.New("Invalid type for containing")
	}
}
