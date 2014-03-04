package state

import "testing"

// Also covers export
func TestMakeScalarContainer(t *testing.T) {
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

// We test "remove with a parent" in other test cases
func TestScalarContainer_Remove(t *testing.T) {
    c, err := MakeScalarContainer("floop")
    if err != nil {
        t.Fatal(err)
    }
    err = c.Remove()
    if err == nil {
        t.Fatal("ScalarContainer.Remove should fail when no parent")
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
    if err != nil { t.Fatal(err) }

    // Break the rules a bit :)
    p, err := MakeScalarContainer("parent")
    if err != nil { t.Fatal(err) }

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
    err = c.Set(0,0)
    if err == nil {
        t.Fatal("ScalarContainer.Set should always fail")
    }
}
