package state

import "errors"

// Even though the SliceContainer represents an array-like value,
// it is easier to internally implement it as a map[uint]Container,
// since this makes it easy to set "the value on the end".
//
// The consequence is a bit more complexity in Export, but that's
// still not so bad, and certainly not as bad as implementing the
// Set method for a []Container slice.
type SliceContainer struct {
	Parent    Container
	ParentKey interface{}
	Value     map[uint]Container
}

func MakeSliceContainer(s []interface{}) (Container, error) {
	c := SliceContainer{
		nil,
		nil,
		make(map[uint]Container),
	}
	for key, value := range s {
		err := c.Set(uint(key), value)
		if err != nil {
			return nil, err
		}
	}
	return &c, nil
}

func (c *SliceContainer) castKey(key interface{}) (uint, error) {
	switch k := key.(type) {
	case uint:
		return k, nil
	case int:
		return uint(k), nil
	default:
		return uint(0), errors.New("Cannot cast key to uint")
	}
}

func (c *SliceContainer) GetChild(key interface{}) (Container, error) {
	key_int, err := c.castKey(key)
	if err != nil {
		return nil, err
	}
	child, ok := c.Value[key_int]
	if !ok {
		return nil, errors.New("Key not present in map")
	}
	return child, nil
}

func (c *SliceContainer) Remove() error {
	if c.Parent == nil {
		return errors.New("No parent")
	}
	return c.Parent.RemoveChild(c.ParentKey)
}

func (c *SliceContainer) RemoveChild(key interface{}) error {
	key_int, err := c.castKey(key)
	if err != nil {
		return err
	}
	delete(c.Value, key_int)
	for key, value := range c.Value {
		if key > key_int {
			delete(c.Value, key)
			c.Value[key-1] = value
		}
	}
	return nil
}

func (c *SliceContainer) SetParentage(p Container, key interface{}) {
	c.Parent = p
	c.ParentKey = key
}

func (c *SliceContainer) Set(key, value interface{}) error {
	key_int, err := c.castKey(key)
	if err != nil {
		return err
	}
	child, err := MakeContainer(value)
	if err != nil {
		return err
	}
	child.SetParentage(c, key_int)
	c.Value[key_int] = child
	return nil
}

func (c *SliceContainer) Export() interface{} {
	// Get total length of result array
	max_key := uint(0)
	for key, _ := range c.Value {
		if key > max_key {
			max_key = key
		}
	}

	// Uninitialized interface{}s are nil
	result := make([]interface{}, max_key+1)
	for key, _ := range result {
		value, ok := c.Value[uint(key)]
		if !ok {
			value, _ = MakeScalarContainer(nil)
		}
		result[key] = value.Export()
	}
	return result
}
