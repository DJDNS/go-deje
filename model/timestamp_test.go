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

func TestTimestampSetContains(t *testing.T) {
	ts := make(TimestampSet)
	A := &Timestamp{
		QHash: "xyz",
		Time:  5,
	}
	B := &Timestamp{
		QHash: "abc",
		Time:  5,
	}

	if ts.Contains(A) {
		t.Fatal("ts doesn't actually contain A")
	}

	ts[A.QHash] = B
	if ts.Contains(A) {
		t.Fatal("ts doesn't actually contain A, but has key for A")
	}

	ts[A.QHash] = A
	if !ts.Contains(A) {
		t.Fatal("ts actually does contain A")
	}
}

func TestTimestampManagerGetBlock(t *testing.T) {
	m := NewTimestampManager()
	block := m.GetBlock(20)

	if len(block) != 0 {
		t.Fatal("block should have been empty")
	}

	ts := Timestamp{
		QHash: "xyz",
		Time:  5,
	}
	m.Register(&ts)

	block = m.GetBlock(20)
	if len(block) != 0 {
		t.Fatal("block should still have been empty")
	}
	block = m.GetBlock(5)
	if len(block) != 1 {
		t.Fatal("block should not have been empty")
	}

	if block[ts.QHash] != &ts {
		t.Fatal("block contents should have contained &ts")
	}

	if len(m.PerBlock) != 2 {
		t.Fatal("Not caching blocks")
	}
}

func TestTimestampManagerRegister(t *testing.T) {
	m := NewTimestampManager()

	A := &Timestamp{
		QHash: "xyz",
		Time:  5,
	}
	B := &Timestamp{
		QHash: "abc",
		Time:  5,
	}

	m.Register(A)
	m.Register(B)

	block := m.GetBlock(5)
	if !(block.Contains(A) && block.Contains(B)) {
		t.Fatal("block was missing timestamps")
	}
}

func TestTimestampManagerUnregister(t *testing.T) {
	m := NewTimestampManager()

	A := &Timestamp{
		QHash: "xyz",
		Time:  5,
	}
	B := &Timestamp{
		QHash: "abc",
		Time:  5,
	}

	m.Register(A)
	m.Register(B)
	m.Unregister(A)

	block := m.GetBlock(5)
	if block.Contains(A) || !block.Contains(B) || len(m.Stamps) != 1 {
		t.Fatal("Failed to unregister A correctly")
	}
}
