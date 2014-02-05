package deje

import "testing"

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

	if !Manageable(A).Eq(A) {
		t.Fatal("A should equal A")
	}
	if Manageable(A).Eq(B) {
		t.Fatal("A should not equal B")
	}
	if !Manageable(A).Eq(C) {
		t.Fatal("A should equal C")
	}
}
