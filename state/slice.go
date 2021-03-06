package state

import "errors"

// Even though the sliceContainer represents an array-like value,
// it is easier to internally implement it as a map[uint]Container,
// since this makes it easy to set "the value on the end".
//
// The consequence is a bit more complexity in Export, but that's
// still not so bad, and certainly not as bad as implementing the
// Set method for a []Container slice.
type sliceContainer struct {
	Value map[uint]Container
}

func makeSliceContainer(s []interface{}) (Container, error) {
	c := sliceContainer{
		make(map[uint]Container),
	}
	for key, value := range s {
		err := c.SetChild(uint(key), value)
		if err != nil {
			return nil, err
		}
	}
	return &c, nil
}

func (c *sliceContainer) castKey(key interface{}) (uint, error) {
	switch k := key.(type) {
	case uint:
		return k, nil
	case int:
		return uint(k), nil
	case float64:
		return uint(k), nil
	default:
		return uint(0), errors.New("Cannot cast key to uint")
	}
}

func (c *sliceContainer) GetChild(key interface{}) (Container, error) {
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

func (c *sliceContainer) SetChild(key, value interface{}) error {
	key_int, err := c.castKey(key)
	if err != nil {
		return err
	}
	child, err := makeContainer(value)
	if err != nil {
		return err
	}
	c.Value[key_int] = child
	return nil
}

func (c *sliceContainer) RemoveChild(key interface{}) error {
	key_int, err := c.castKey(key)
	if err != nil {
		return err
	}
	delete(c.Value, key_int)

	// Get max index
	var max uint
	for k := range c.Value {
		if k > max {
			max = k
		}
	}

	// Iterate from key+1 to max, decrementing all indexes after key IN ORDER
	for k := key_int + 1; k <= max; k++ {
		value, ok := c.Value[k]
		if ok {
			delete(c.Value, k)
			c.Value[k-1] = value
		}
	}
	return nil
}

func (c *sliceContainer) Export() interface{} {
	// Fast - and correct - special case.
	// The general purpose code is actually broken for empty c.Value.
	if len(c.Value) == 0 {
		return make([]interface{}, 0)
	}

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
			value, _ = makeScalarContainer(nil)
		}
		result[key] = value.Export()
	}
	return result
}
