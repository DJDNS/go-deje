package state

import (
	"reflect"
	"testing"
)

// We test with valid children all over the place.
// Let's break things.
func TestMakeSliceContainer_InvalidChildren(t *testing.T) {
	original := []interface{}{
		make(chan int),
	}
	_, err := makeSliceContainer(original)
	if err == nil {
		t.Fatal("makeSliceContainer should fail if it can't contain children")
	}
}

func TestSliceContainer_GetChild(t *testing.T) {
	original := []interface{}{
		"hello",
		"world",
	}
	c, err := makeSliceContainer(original)
	if err != nil {
		t.Fatal(err)
	}
	_, err = c.GetChild(true)
	if err == nil {
		t.Fatal("Non-uint keys should fail")
	}
	_, err = c.GetChild(uint(8))
	if err == nil {
		t.Fatal("Should fail getting from unset key")
	}

	child, err := c.GetChild(uint(0))
	if err != nil {
		t.Fatal(err)
	}
	if child.Export() != "hello" {
		t.Fatalf("Should have got hello child, got %#v", child)
	}

	child, err = c.GetChild(uint(1))
	if err != nil {
		t.Fatal(err)
	}
	if child.Export() != "world" {
		t.Fatalf("Should have got world child, got %#v", child)
	}
}

func TestSliceContainer_SetChild(t *testing.T) {
	c, err := makeSliceContainer([]interface{}{})
	if err != nil {
		t.Fatal(err)
	}
	err = c.SetChild("hi", 0)
	if err == nil {
		t.Fatal("sliceContainer.SetChild with non-uint key should always fail")
	}
	err = c.SetChild(uint(9), make(chan int))
	if err == nil {
		t.Fatal("sliceContainer.SetChild with non-containable value should fail")
	}
	err = c.SetChild(uint(5), 89)
	if err != nil {
		t.Fatal(err)
	}
	expected := []interface{}{
		nil, nil, nil, nil, nil, 89,
	}
	if !reflect.DeepEqual(c.Export(), expected) {
		t.Fatal("Expected %#v, got %#v", expected, c.Export())
	}
}

func TestSliceContainer_RemoveChild(t *testing.T) {
	original := []interface{}{
		"hello", "crazy", "world",
	}
	c, err := makeSliceContainer(original)
	if err != nil {
		t.Fatal(err)
	}
	err = c.RemoveChild("flub")
	if err == nil {
		t.Fatal("sliceContainer.RemoveChild should fail for bad key type")
	}
	err = c.RemoveChild(uint(90))
	if err != nil {
		t.Fatal(err) // Should NOT fail for unset keys
	}
	err = c.RemoveChild(uint(0))
	if err != nil {
		t.Fatal(err)
	}
	expected := []interface{}{
		"crazy", "world",
	}
	if !reflect.DeepEqual(c.Export(), expected) {
		t.Fatalf("Expected %#v, got %#v", expected, c.Export())
	}
}

func TestSliceContainerExport(t *testing.T) {
	demo_map := []interface{}{
		"Hello", "world",
		"Nine", 9,
		[]interface{}{
			true, false, nil,
		},
		map[string]interface{}{
			"deep": "path",
		},
	}
	container, err := makeSliceContainer(demo_map)
	if err != nil {
		t.Fatal(err)
	}
	exported := container.Export()
	if !reflect.DeepEqual(exported, demo_map) {
		t.Fatalf("Expected %v, got %v", demo_map, exported)
	}
}
