package model

import "testing"

func TestMS_Contains(t *testing.T) {
	A := Quorum{"A", nil}
	B := Quorum{"B", nil}
	C := Quorum{"C", nil}

	ms := make(ManageableSet)

	ms[A.GetKey()] = A // A registered correctly
	ms[B.GetKey()] = A // B registered with wrong value
	// C not registered at all

	if !ms.Contains(A) {
		t.Fatal("ms should contain A")
	}
	if ms.Contains(B) {
		t.Fatal("ms should not contain B")
	}
	if ms.Contains(C) {
		t.Fatal("ms should not contain C")
	}
}
