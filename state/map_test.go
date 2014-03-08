package state

import (
	"reflect"
	"testing"
)

// We test with valid children all over the place.
// Let's break things.
func TestMakeMapContainer_InvalidChildren(t *testing.T) {
	original := map[string]interface{}{
		"key": make(chan int),
	}
	_, err := MakeMapContainer(original)
	if err == nil {
		t.Fatal("MakeMapContainer should fail if it can't contain children")
	}
}

func TestMapContainer_GetChild(t *testing.T) {
	original := map[string]interface{}{
		"hello": "world",
		"recursive": map[string]interface{}{
			"deep": "path",
		},
	}
	c, err := MakeMapContainer(original)
	if err != nil {
		t.Fatal(err)
	}
	_, err = c.GetChild(true)
	if err == nil {
		t.Fatal("Non-string keys should fail")
	}
	_, err = c.GetChild("floop")
	if err == nil {
		t.Fatal("Should fail getting from unset key")
	}

	child, err := c.GetChild("hello")
	if err != nil {
		t.Fatal(err)
	}
	if child.Export() != "world" {
		t.Fatalf("Should have got string child, got %#v", child)
	}

	child, err = c.GetChild("recursive")
	if err != nil {
		t.Fatal(err)
	}
	expected := map[string]interface{}{"deep": "path"}
	if !reflect.DeepEqual(child.Export(), expected) {
		t.Fatalf("Should have got map child, got %#v", child)
	}
}

func TestMapContainer_SetChild(t *testing.T) {
	c, err := MakeMapContainer(map[string]interface{}{})
	if err != nil {
		t.Fatal(err)
	}
	err = c.SetChild(0, 0)
	if err == nil {
		t.Fatal("MapContainer.SetChild with non-str key should always fail")
	}
	err = c.SetChild("some_key", make(chan int))
	if err == nil {
		t.Fatal("MapContainer.SetChild with non-containable value should fail")
	}
	err = c.SetChild("some_key", 0)
	if err != nil {
		t.Fatal(err)
	}
	expected := map[string]interface{}{
		"some_key": 0,
	}
	if !reflect.DeepEqual(c.Export(), expected) {
		t.Fatal("Expected %#v, got %#v", expected, c.Export())
	}
}

func TestMapContainer_RemoveChild(t *testing.T) {
	original := map[string]interface{}{
		"hello": "world",
		"recursive": map[string]interface{}{
			"deep": "path",
		},
	}
	c, err := MakeMapContainer(original)
	if err != nil {
		t.Fatal(err)
	}
	err = c.RemoveChild(0)
	if err == nil {
		t.Fatal("MapContainer.RemoveChild should fail for bad key type")
	}
	err = c.RemoveChild("floop")
	if err != nil {
		t.Fatal(err) // Should NOT fail for unset keys
	}
	err = c.RemoveChild("hello")
	if err != nil {
		t.Fatal(err)
	}
	expected := map[string]interface{}{
		"recursive": map[string]interface{}{
			"deep": "path",
		},
	}
	if !reflect.DeepEqual(c.Export(), expected) {
		t.Fatal("Item should be removed")
	}
}

func TestMapContainerExport(t *testing.T) {
	demo_map := map[string]interface{}{
		"Hello": "world",
		"Nine":  9,
		"recursive": map[string]interface{}{
			"deep": "path",
		},
	}
	container, err := MakeMapContainer(demo_map)
	if err != nil {
		t.Fatal(err)
	}
	exported := container.Export()
	if !reflect.DeepEqual(exported, demo_map) {
		t.Fatalf("Expected %v, got %v", demo_map, exported)
	}
}
