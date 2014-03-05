package state

import (
	"reflect"
	"testing"
)

func TestScalarContainer_GetChild(t *testing.T) {
	c, err := MakeScalarContainer("floop")
	if err != nil {
		t.Fatal(err)
	}
	_, err = c.GetChild(0)
	if err == nil {
		t.Fatal("ScalarContainer.GetChild should always fail")
	}
}

func TestScalarContainer_Remove_NoParent(t *testing.T) {
	c, err := MakeScalarContainer("floop")
	if err != nil {
		t.Fatal(err)
	}
	err = c.Remove()
	if err == nil {
		t.Fatal("ScalarContainer.Remove should fail when no parent")
	}
}

func TestScalarContainer_Remove_WithParent(t *testing.T) {
	map_original := map[string]interface{}{
		"hello": "world",
		"this":  "that",
	}
	c, err := MakeMapContainer(map_original)
	if err != nil {
		t.Fatal(err)
	}
	child, err := c.GetChild("hello")
	if err != nil {
		t.Fatal(err)
	}
	err = child.Remove()
	if err != nil {
		t.Fatal(err)
	}
	expected := map[string]interface{}{
		"this": "that",
	}
	if !reflect.DeepEqual(c.Export(), expected) {
		t.Fatal("Item should be removed")
	}
}

func TestScalarContainer_RemoveChild(t *testing.T) {
	c, err := MakeScalarContainer("floop")
	if err != nil {
		t.Fatal(err)
	}
	err = c.RemoveChild(0)
	if err == nil {
		t.Fatal("ScalarContainer.RemoveChild should always fail")
	}
}

func TestScalarContainer_SetParentage(t *testing.T) {
	c, err := MakeScalarContainer("floop")
	if err != nil {
		t.Fatal(err)
	}

	// Break the rules a bit :)
	p, err := MakeScalarContainer("parent")
	if err != nil {
		t.Fatal(err)
	}

	c.SetParentage(p, true)
	if c.(*ScalarContainer).Parent != p {
		t.Fatal("c.Parent should equal p")
	}
	if c.(*ScalarContainer).ParentKey != true {
		t.Fatal("c.ParentKey should equal true")
	}
}

func TestScalarContainer_Set(t *testing.T) {
	c, err := MakeScalarContainer("floop")
	if err != nil {
		t.Fatal(err)
	}
	err = c.Set(0, 0)
	if err == nil {
		t.Fatal("ScalarContainer.Set should always fail")
	}
}

func TestScalarContainer_Export(t *testing.T) {
	scalars := []interface{}{
		"hello",
		nil,
		80,
		true,
		false,
	}

	for _, scalar := range scalars {
		c, err := MakeScalarContainer(scalar)
		if err != nil {
			t.Fatal(err)
		}
		if c.Export() != scalar {
			t.Fatalf("Expected %v, got %v", scalar, c.Export())
		}
	}
}
