package model

import "testing"

func TestQuorum_GetKey(t *testing.T) {
	A := NewQuorum("hello")
	B := NewQuorum("hello")

	if A.GetKey() != B.GetKey() {
		t.Fatal("A should equal B")
	}

	A.Signatures["ident"] = "sig"
	if A.GetKey() == B.GetKey() {
		t.Fatal("A should not equal B")
	}
}

func TestQuorum_GetGroupKey(t *testing.T) {
	A := NewQuorum("hello")
	B := NewQuorum("hello")

	if A.GetGroupKey() != B.GetGroupKey() {
		t.Fatal("A should equal B")
	}

	// Signature content should not change group key
	A.Signatures["ident"] = "sig"
	if A.GetGroupKey() != B.GetGroupKey() {
		t.Fatal("A should equal B")
	}

	// Different EventHash should alter group key
	A.EventHash = "world"
	if A.GetGroupKey() == B.GetGroupKey() {
		t.Fatal("A should not equal B")
	}
}

func TestQuorum_Eq(t *testing.T) {
	Q1 := NewQuorum("hello")
	Q2 := NewQuorum("hello")
	E := NewEvent("handler name")

	if !Q1.Eq(Q2) {
		t.Fatal("Q1 should equal Q2")
	}

	Q1.Signatures["ident"] = "sig"
	if Q1.Eq(Q2) {
		t.Fatal("Q1 should not equal Q2")
	}

	if Q1.Eq(E) {
		t.Fatal("Q1 should not equal E")
	}
}
