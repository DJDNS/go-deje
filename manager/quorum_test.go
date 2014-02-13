package model

import "testing"

func TestQuorumManager_Register(t *testing.T) {
	m := NewQuorumManager()
	q := NewQuorum("evhash")

	if m.Contains(q) {
		t.Fatal("m should not contain q yet")
	}
	m.Register(q)
	if !m.Contains(q) {
		t.Fatal("m should contain q")
	}
}

func TestQuorumManager_Unregister(t *testing.T) {
	m := NewQuorumManager()
	q := NewQuorum("evhash")
	m.Register(q)

	if !m.Contains(q) {
		t.Fatal("m should contain q")
	}
	m.Unregister(q)
	if m.Contains(q) {
		t.Fatal("m should not contain q anymore")
	}

	// Should be idempotent
	m.Unregister(q)
	if m.Contains(q) {
		t.Fatal("m should not contain q anymore")
	}
}
