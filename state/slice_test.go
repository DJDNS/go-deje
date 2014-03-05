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
	_, err := MakeSliceContainer(original)
	if err == nil {
		t.Fatal("MakeSliceContainer should fail if it can't contain children")
	}
}

func TestSliceContainer_GetChild(t *testing.T) {
	original := []interface{}{
		"hello",
		"world",
	}
	c, err := MakeSliceContainer(original)
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

func TestSliceContainer_Remove_NoParent(t *testing.T) {
	original := []interface{}{
		"hello", "world",
	}
	c, err := MakeSliceContainer(original)
	if err != nil {
		t.Fatal(err)
	}
	err = c.Remove()
	if err == nil {
		t.Fatal("SliceContainer.Remove should fail when no parent")
	}
}

func TestSliceContainer_Remove_WithParent(t *testing.T) {
	original := []interface{}{
		"hello",
		[]interface{}{"sublist", "items"},
		"world",
	}
	c, err := MakeSliceContainer(original)
	if err != nil {
		t.Fatal(err)
	}
	child, err := c.GetChild(uint(1))
	if err != nil {
		t.Fatal(err)
	}
	err = child.Remove()
	if err != nil {
		t.Fatal(err)
	}
	expected := []interface{}{
		"hello", "world",
	}
	if !reflect.DeepEqual(c.Export(), expected) {
		t.Fatalf("Expected %#v, got %#v", expected, c.Export())
	}
}

func TestSliceContainer_RemoveChild(t *testing.T) {
	original := []interface{}{
		"hello", "crazy", "world",
	}
	c, err := MakeSliceContainer(original)
	if err != nil {
		t.Fatal(err)
	}
	err = c.RemoveChild("flub")
	if err == nil {
		t.Fatal("SliceContainer.RemoveChild should fail for bad key type")
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

func TestSliceContainer_SetParentage(t *testing.T) {
	c, err := MakeSliceContainer([]interface{}{})
	if err != nil {
		t.Fatal(err)
	}

	// Break the rules a bit :)
	p, err := MakeScalarContainer("parent")
	if err != nil {
		t.Fatal(err)
	}

	c.SetParentage(p, true)
	if c.(*SliceContainer).Parent != p {
		t.Fatal("c.Parent should equal p")
	}
	if c.(*SliceContainer).ParentKey != true {
		t.Fatal("c.ParentKey should equal true")
	}
}

func TestSliceContainer_Set(t *testing.T) {
	c, err := MakeSliceContainer([]interface{}{})
	if err != nil {
		t.Fatal(err)
	}
	err = c.Set("hi", 0)
	if err == nil {
		t.Fatal("SliceContainer.Set with non-uint key should always fail")
	}
	err = c.Set(uint(9), make(chan int))
	if err == nil {
		t.Fatal("SliceContainer.Set with non-containable value should fail")
	}
	err = c.Set(uint(5), 89)
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
	container, err := MakeSliceContainer(demo_map)
	if err != nil {
		t.Fatal(err)
	}
	exported := container.Export()
	if !reflect.DeepEqual(exported, demo_map) {
		t.Fatalf("Expected %v, got %v", demo_map, exported)
	}
}
