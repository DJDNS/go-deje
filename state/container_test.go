package state

import "testing"

// We cover valid maps elsewhere. Let's be bad.
func TestMakeContainer_InvalidMap(t *testing.T) {
	original := make(map[int]string)
	_, err := MakeContainer(original)
	if err == nil {
		t.Fatal("MakeContainer should fail for invalid map type")
	}

	// Should be map[string]interface{}
	other_invalid := make(map[string]int)
	_, err = MakeContainer(other_invalid)
	if err == nil {
		t.Fatal("MakeContainer should fail for invalid map type")
	}
}
