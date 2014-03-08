package state

import "testing"

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

func TestScalarContainer_SetChild(t *testing.T) {
	c, err := MakeScalarContainer("floop")
	if err != nil {
		t.Fatal(err)
	}
	err = c.SetChild(0, 0)
	if err == nil {
		t.Fatal("ScalarContainer.SetChild should always fail")
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
