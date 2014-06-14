package state

import "errors"

type scalarContainer struct {
	Value interface{}
}

func makeScalarContainer(value interface{}) (Container, error) {
	return &scalarContainer{
		value,
	}, nil
}

func (c *scalarContainer) SetChild(key, value interface{}) error {
	return errors.New("Scalars do not have children")
}

func (c *scalarContainer) GetChild(key interface{}) (Container, error) {
	return nil, errors.New("Scalars do not have children")
}

func (c *scalarContainer) RemoveChild(key interface{}) error {
	return errors.New("Scalars do not have children")
}

func (c *scalarContainer) Export() interface{} {
	return c.Value
}
