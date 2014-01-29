package deje

import "testing"

func TestTimestampComparison(t *testing.T) {
	A := Timestamp{
		SyncHash:    "xyz",
		BlockHeight: 5,
	}
	B := Timestamp{
		SyncHash:    "abc",
		BlockHeight: 4,
	}

	if !(B.WasBefore(A) && !A.WasBefore(B)) {
		t.Fatal("A < B")
	}

	// Tied on blockheight, but B still has lower hash
	B.BlockHeight = 5
	if !(B.WasBefore(A) && !A.WasBefore(B)) {
		t.Fatal("A < B")
	}

	// B happens a block after A
	B.BlockHeight = 6
	if !(A.WasBefore(B) && !B.WasBefore(A)) {
		t.Fatal("A > B")
	}
}
