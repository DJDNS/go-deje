package state

import "errors"

type MapContainer struct {
	Parent    Container
	ParentKey interface{}
	Value     map[string]Container
}

func MakeMapContainer(m map[string]interface{}) (Container, error) {
	c := MapContainer{
		nil,
		nil,
		make(map[string]Container),
	}
	for key, value := range m {
		err := c.Set(key, value)
		if err != nil {
			return nil, err
		}
	}
	return &c, nil
}

func (c *MapContainer) GetChild(key interface{}) (Container, error) {
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

func (c *MapContainer) Remove() error {
	if c.Parent == nil {
		return errors.New("No parent")
	}
	return c.Parent.RemoveChild(c.ParentKey)
}

func (c *MapContainer) RemoveChild(key interface{}) error {
	key_str, ok := key.(string)
	if !ok {
		return errors.New("Key was not string type")
	}
	delete(c.Value, key_str)
	return nil
}

func (c *MapContainer) SetParentage(p Container, key interface{}) {
	c.Parent = p
	c.ParentKey = key
}

func (c *MapContainer) Set(key, value interface{}) error {
	key_str, ok := key.(string)
	if !ok {
		return errors.New("Key was not string type")
	}
	child, err := MakeContainer(value)
	if err != nil {
		return err
	}
	child.SetParentage(c, key_str)
	c.Value[key_str] = child
	return nil
}

func (c *MapContainer) Export() interface{} {
	result := make(map[string]interface{})
	for key, value := range c.Value {
		result[key] = value.Export()
	}
	return result
}
