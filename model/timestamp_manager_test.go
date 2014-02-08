package model

import "testing"

func TestTimestampManager_Register(t *testing.T) {
	m := NewTimestampManager()
	ts := Timestamp{
		QHash: "Bubba",
		Time:  25,
	}

	if m.Contains(ts) {
		t.Fatal("m should not contain ts yet")
	}
	m.Register(ts)
	if !m.Contains(ts) {
		t.Fatal("m should contain ts")
	}
}

func TestTimestampManager_Unregister(t *testing.T) {
	m := NewTimestampManager()
	ts := Timestamp{
		QHash: "Bubba",
		Time:  25,
	}
	m.Register(ts)

	if !m.Contains(ts) {
		t.Fatal("m should contain ts")
	}
	m.Unregister(ts)
	if m.Contains(ts) {
		t.Fatal("m should not contain ts anymore")
	}

	// Should be idempotent
	m.Unregister(ts)
	if m.Contains(ts) {
		t.Fatal("m should not contain ts anymore")
	}
}
