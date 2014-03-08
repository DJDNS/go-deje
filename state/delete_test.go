package state

import "testing"

func TestDeletePrimitive_Apply_Root(t *testing.T) {
	ds := NewDocumentState()
	primitive := DeletePrimitive{
		Path: []interface{}{},
	}
	err := primitive.Apply(ds)
	if err == nil {
		t.Fatal("DeletePrimitive.Apply should always fail for root")
	}
}

func TestDeletePrimitive_Apply_WithPath(t *testing.T) {
	ds := NewDocumentState()
	primitive := DeletePrimitive{
		Path: []interface{}{"this", "that"},
	}
	err := primitive.Apply(ds)
	if err == nil {
		t.Fatal("DeletePrimitive.Apply should fail for bad traversals")
	}
	primitive.Path = []interface{}{"this"}
	err = primitive.Apply(ds)
	if err != nil {
		t.Fatal(err) // Should not fail - deletion is idempotent
	}

	// Confirm that actual deletion happens
	setter := SetPrimitive{
		Path: []interface{}{},
		Value: map[string]interface{}{
			"hello": "world",
			"other": "stuff",
			"deep": map[string]interface{}{
				"this":   "that",
				"things": []interface{}{true, false, nil},
			},
		},
	}
	err = setter.Apply(ds)
	if err != nil {
		t.Fatal(err)
	}

	tests := []primitiveApplyTest{
		primitiveApplyTest{
			&DeletePrimitive{[]interface{}{"this"}},
			map[string]interface{}{
				"hello": "world",
				"other": "stuff",
				"deep": map[string]interface{}{
					"this":   "that",
					"things": []interface{}{true, false, nil},
				},
			},
			"Deleting nonexistent key should have no effect",
		},
		primitiveApplyTest{
			&DeletePrimitive{[]interface{}{"other"}},
			map[string]interface{}{
				"hello": "world",
				"deep": map[string]interface{}{
					"this":   "that",
					"things": []interface{}{true, false, nil},
				},
			},
			"Should have removed 'other' item from DS",
		},
		primitiveApplyTest{
			&DeletePrimitive{[]interface{}{"deep", "things", 1}},
			map[string]interface{}{
				"hello": "world",
				"deep": map[string]interface{}{
					"this":   "that",
					"things": []interface{}{true, nil},
				},
			},
			"Should have removed deep-traversal item from DS",
		},
	}
	for _, test := range tests {
		test.Run(t, ds)
	}
}

func TestDeletePrimitive_Reverse(t *testing.T) {
	tests := []primitiveReverseTest{
		primitiveReverseTest{
			Original: map[string]interface{}{},
			Change: &DeletePrimitive{
				Path: []interface{}{"missing"},
			},
			FailureMsg: "Reversal on root change",
		},
		primitiveReverseTest{
			Original: map[string]interface{}{
				"existing": "stuff",
			},
			Change: &DeletePrimitive{
				Path: []interface{}{"existing"},
			},
			FailureMsg: "Reversal on root restores old root",
		},
		primitiveReverseTest{
			Original: map[string]interface{}{
				"existing": "stuff",
				"hello":    "world",
			},
			Change: &DeletePrimitive{
				Path: []interface{}{"hello"},
			},
			FailureMsg: "Reversal on root does not affect existing",
		},
	}

	for _, test := range tests {
		test.Run(t)
	}
}
