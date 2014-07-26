package state

import (
	"encoding/json"
	"reflect"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
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
func TestMakeSliceContainer_NoChildren(t *testing.T) {
	original := []interface{}{}
	sc, err := makeSliceContainer(original)
	if err != nil {
		t.Fatal(err)
	}
	assert.Equal(t, original, sc.Export())
}

func TestSliceContainer_castKey_JsonNumber(t *testing.T) {
	container, err := makeSliceContainer([]interface{}{})
	if err != nil {
		t.Fatal(err)
	}

	var array []interface{}
	if err = json.Unmarshal([]byte("[84, 12.1]"), &array); err != nil {
		t.Fatal(err)
	}
	for i := range array {
		number, err := container.(*sliceContainer).castKey(array[i])
		if err != nil {
			t.Fatal(err)
		}
		array[i] = number
	}

	expected_array := []interface{}{uint(84), uint(12)}
	if !reflect.DeepEqual(array, expected_array) {
		t.Fatalf("Expected %#v, got %#v", expected_array, array)
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

// A bit of statistical surety that we are decrementing in the correct order.
// https://github.com/DJDNS/go-deje/issues/30
func TestSliceContainer_RemoveChild_Many(t *testing.T) {
	text := "The quick brown fox jumped over the lazy dog."
	expected := "The quick brown ox jumped over the lazy dog."
	letters := strings.Split(text, "")

	interface_letters := make([]interface{}, len(letters))
	for i, v := range letters {
		interface_letters[i] = v
	}

	// Turn letters slice into container, remove letter 16 (f)
	c, err := makeSliceContainer(interface_letters)
	if err != nil {
		t.Fatal(err)
	}
	if err := c.RemoveChild(uint(16)); err != nil {
		t.Fatal(err)
	}

	// Reassemble into single string
	export := c.Export().([]interface{})
	strings_export := make([]string, len(export))
	for i, v := range export {
		str, ok := v.(string)
		if !ok {
			t.Fatalf("Bad values in array: %#v", export)
		}
		strings_export[i] = str
	}
	result := strings.Join(strings_export, "")

	// Test for clobbering
	assert.Equal(t, expected, result)
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
