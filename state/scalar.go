package state

import "errors"

type ScalarContainer struct {
	Value interface{}
}

func MakeScalarContainer(value interface{}) (Container, error) {
	return &ScalarContainer{
		value,
	}, nil
}

func (c *ScalarContainer) SetChild(key, value interface{}) error {
	return errors.New("Scalars do not have children")
}

func (c *ScalarContainer) GetChild(key interface{}) (Container, error) {
	return nil, errors.New("Scalars do not have children")
}

func (c *ScalarContainer) RemoveChild(key interface{}) error {
	return errors.New("Scalars do not have children")
}

func (c *ScalarContainer) Export() interface{} {
	return c.Value
}
