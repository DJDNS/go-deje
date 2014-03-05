package state

import "errors"

type ScalarContainer struct {
	Parent    Container
	ParentKey interface{}
	Value     interface{}
}

func MakeScalarContainer(value interface{}) (Container, error) {
	return &ScalarContainer{
		nil,
		nil,
		value,
	}, nil
}

func (c *ScalarContainer) GetChild(key interface{}) (Container, error) {
	return nil, errors.New("Scalars do not have children")
}

func (c *ScalarContainer) Remove() error {
	if c.Parent == nil {
		return errors.New("No parent")
	}
	return c.Parent.RemoveChild(c.ParentKey)
}

func (c *ScalarContainer) RemoveChild(key interface{}) error {
	return errors.New("Scalars do not have children")
}

func (c *ScalarContainer) SetParentage(p Container, key interface{}) {
	c.Parent = p
	c.ParentKey = key
}

func (c *ScalarContainer) Set(key, value interface{}) error {
	return errors.New("Scalars do not have children")
}

func (c *ScalarContainer) Export() interface{} {
	return c.Value
}
