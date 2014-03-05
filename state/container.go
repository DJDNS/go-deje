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
	switch reflect.TypeOf(value).Kind() {
	case reflect.Map:
		as_map, ok := value.(map[string]interface{})
		if !ok {
			return nil, errors.New("Cannot cast map to map[string]interface{}")
		}
		return MakeMapContainer(as_map)
	case reflect.Bool, reflect.Int, reflect.Uint, reflect.String:
		return MakeScalarContainer(value)
	default:
		return nil, errors.New("Invalid type for containing")
	}
}
