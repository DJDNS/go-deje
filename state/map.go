package state

import "errors"

type mapContainer struct {
	Value map[string]Container
}

func makeMapContainer(m map[string]interface{}) (Container, error) {
	c := mapContainer{
		make(map[string]Container),
	}
	for key, value := range m {
		err := c.SetChild(key, value)
		if err != nil {
			return nil, err
		}
	}
	return &c, nil
}

func (c *mapContainer) GetChild(key interface{}) (Container, error) {
	key_str, ok := key.(string)
	if !ok {
		return nil, errors.New("Key was not string type")
	}
	child, ok := c.Value[key_str]
	if !ok {
		return nil, errors.New("Key not present in map")
	}
	return child, nil
}

func (c *mapContainer) SetChild(key, value interface{}) error {
	key_str, ok := key.(string)
	if !ok {
		return errors.New("Key was not string type")
	}
	child, err := makeContainer(value)
	if err != nil {
		return err
	}
	c.Value[key_str] = child
	return nil
}

func (c *mapContainer) RemoveChild(key interface{}) error {
	key_str, ok := key.(string)
	if !ok {
		return errors.New("Key was not string type")
	}
	delete(c.Value, key_str)
	return nil
}

func (c *mapContainer) Export() interface{} {
	result := make(map[string]interface{})
	for key, value := range c.Value {
		result[key] = value.Export()
	}
	return result
}
