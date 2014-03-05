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

func TestMapContainer_Remove_NoParent(t *testing.T) {
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
	err = c.Remove()
	if err == nil {
		t.Fatal("MapContainer.Remove should fail when no parent")
	}
}

func TestMapContainer_Remove_WithParent(t *testing.T) {
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
	child, err := c.GetChild("recursive")
	if err != nil {
		t.Fatal(err)
	}
	err = child.Remove()
	if err != nil {
		t.Fatal(err)
	}
	expected := map[string]interface{}{
		"hello": "world",
	}
	if !reflect.DeepEqual(c.Export(), expected) {
		t.Fatal("Item should be removed")
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

func TestMapContainer_SetParentage(t *testing.T) {
	c, err := MakeMapContainer(map[string]interface{}{})
	if err != nil {
		t.Fatal(err)
	}

	// Break the rules a bit :)
	p, err := MakeScalarContainer("parent")
	if err != nil {
		t.Fatal(err)
	}

	c.SetParentage(p, true)
	if c.(*MapContainer).Parent != p {
		t.Fatal("c.Parent should equal p")
	}
	if c.(*MapContainer).ParentKey != true {
		t.Fatal("c.ParentKey should equal true")
	}
}

func TestMapContainer_Set(t *testing.T) {
	c, err := MakeMapContainer(map[string]interface{}{})
	if err != nil {
		t.Fatal(err)
	}
	err = c.Set(0, 0)
	if err == nil {
		t.Fatal("MapContainer.Set with non-str key should always fail")
	}
	err = c.Set("some_key", make(chan int))
	if err == nil {
		t.Fatal("MapContainer.Set with non-containable value should fail")
	}
	err = c.Set("some_key", 0)
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
