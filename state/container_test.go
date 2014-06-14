package state

import "testing"

// We cover valid maps elsewhere. Let's be bad.
func TestMakeContainer_InvalidMap(t *testing.T) {
	original := make(map[int]string)
	_, err := makeContainer(original)
	if err == nil {
		t.Fatal("makeContainer should fail for invalid map type")
	}

	// Should be map[string]interface{}
	other_invalid := make(map[string]int)
	_, err = makeContainer(other_invalid)
	if err == nil {
		t.Fatal("makeContainer should fail for invalid map type")
	}
}

// Similarly, test for non-interface{} slices
func TestMakeContainer_InvalidSlice(t *testing.T) {
	original := make([]string, 0)
	_, err := makeContainer(original)
	if err == nil {
		t.Fatal("makeContainer should fail for invalid slice type")
	}
}

func TestTraverse(t *testing.T) {
	original := map[string]interface{}{
		"deep": []interface{}{
			"chain", "of", map[string]interface{}{
				"stored": "stuff",
			},
		},
	}

	root, err := makeContainer(original)
	if err != nil {
		t.Fatal(err)
	}

	child, err := Traverse(root, []interface{}{"bad key"})
	if err == nil {
		t.Fatal("Traverse should fail on bad key")
	}
	child, err = Traverse(root, []interface{}{0})
	if err == nil {
		t.Fatal("Traverse should fail on bad key type")
	}
	child, err = Traverse(root, []interface{}{"deep", 0, "baloney"})
	if err == nil {
		t.Fatal("Traverse should fail on bad traversal")
	}
	child, err = Traverse(root, []interface{}{"deep", 0})
	if err != nil {
		t.Fatal(err)
	}
	if child.Export() != "chain" {
		t.Fatal("Expected chain string, got %v", child.Export())
	}
	child, err = Traverse(root, []interface{}{"deep", 2, "stored"})
	if err != nil {
		t.Fatal(err)
	}
	if child.Export() != "stuff" {
		t.Fatal("Expected stuff string, got %v", child.Export())
	}
}
