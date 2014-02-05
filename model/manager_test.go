package model

import "testing"

func TestObjectManagerRegister(t *testing.T) {
	m := NewObjectManager()

	A := Timestamp{
		QHash: "xyz",
		Time:  5,
	}
	B := Timestamp{
		QHash: "abc",
		Time:  5,
	}

	m.Register(A)
	m.Register(B)

	group := m.GetGroup("5")
	if !(group.Contains(A) && group.Contains(B)) {
		t.Fatal("group was missing timestamps")
	}
}

func TestObjectManagerUnregister(t *testing.T) {
	m := NewObjectManager()

	A := Timestamp{
		QHash: "xyz",
		Time:  5,
	}
	B := Timestamp{
		QHash: "abc",
		Time:  5,
	}

	m.Register(A)
	m.Register(B)
	m.Unregister(A)

	group := m.GetGroup("5")
	if group.Contains(A) || !group.Contains(B) || m.Length() != 1 {
		t.Fatal("Failed to unregister A correctly")
	}
}

func TestManagableSetContains(t *testing.T) {
	ts := make(ManageableSet)
	A := Timestamp{
		QHash: "xyz",
		Time:  5,
	}
	B := Timestamp{
		QHash: "abc",
		Time:  5,
	}

	if ts.Contains(A) {
		t.Fatal("ts shouldn't contain A")
	}

	ts[A.QHash] = B
	if ts.Contains(A) {
		t.Fatal("ts shouldn't contain A, but has key for A")
	}

	ts[A.QHash] = A
	if !ts.Contains(A) {
		t.Fatal("ts should contain A")
	}
}

func TestObjectManagerContains(t *testing.T) {
	m := NewObjectManager()
	A := Timestamp{
		QHash: "xyz",
		Time:  5,
	}
	B := Timestamp{
		QHash: "abc",
		Time:  5,
	}

	m.Register(A)

	if !m.Contains(A) {
		t.Fatal("m should contain A")
	}
	if m.Contains(B) {
		t.Fatal("m should not contain B")
	}
}

func TestObjectManagerGetGroup(t *testing.T) {
	m := NewObjectManager()
	group := m.GetGroup("5")

	if len(group) != 0 {
		t.Fatal("group should have been empty")
	}

	ts := Timestamp{
		QHash: "xyz",
		Time:  5,
	}
	m.Register(&ts)

	group = m.GetGroup("20")
	if len(group) != 0 {
		t.Fatal("group should still have been empty")
	}
	group = m.GetGroup("5")
	if len(group) != 1 {
		t.Fatal("group should not have been empty")
	}

	if group[ts.QHash] != &ts {
		t.Fatal("group contents should have contained &ts")
	}

	if len(m.by_group) != 2 {
		t.Fatalf(
			"Not caching groups - expected 2 groups, got %d: %v",
			len(m.by_group),
			m.by_group,
		)
	}
}
