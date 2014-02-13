package model

import "testing"

func TestTimestampGetKey(t *testing.T) {
	ts := Timestamp{
		QHash: "xyz",
		Time:  230,
	}

	key := ts.GetKey()
	expected := ts.QHash

	if key != expected {
		t.Fatalf("Expected key %v, got %v", expected, key)
	}
}

func TestTimestampGetGroupKey(t *testing.T) {
	ts := Timestamp{
		QHash: "xyz",
		Time:  230,
	}

	key := ts.GetGroupKey()
	expected := "230"

	if key != expected {
		t.Fatalf("Expected key %v, got %v", expected, key)
	}
}

func TestTimestampEq(t *testing.T) {
	A := Timestamp{
		QHash: "xyz",
		Time:  5,
	}
	B := Timestamp{
		QHash: "abc",
		Time:  4,
	}
	C := Timestamp{
		QHash: "xyz",
		Time:  5,
	}
	Q := NewQuorum("whatever")

	if !Manageable(A).Eq(A) {
		t.Fatal("A should equal A")
	}
	if Manageable(A).Eq(B) {
		t.Fatal("A should not equal B")
	}
	if !Manageable(A).Eq(C) {
		t.Fatal("A should equal C")
	}
	if Manageable(A).Eq(Q) {
		t.Fatal("A should not equal Q")
	}
}

func TestTimestampComparison(t *testing.T) {
	A := Timestamp{
		QHash: "xyz",
		Time:  5,
	}
	B := Timestamp{
		QHash: "abc",
		Time:  4,
	}

	if !(B.WasBefore(A) && !A.WasBefore(B)) {
		t.Fatal("A < B")
	}

	// Tied on blockheight, but B still has lower hash
	B.Time = 5
	if !(B.WasBefore(A) && !A.WasBefore(B)) {
		t.Fatal("A < B")
	}

	// B happens a block after A
	B.Time = 6
	if !(A.WasBefore(B) && !B.WasBefore(A)) {
		t.Fatal("A > B")
	}
}
