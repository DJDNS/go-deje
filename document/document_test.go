package document

import (
	"testing"
)

func TestNewDocument(t *testing.T) {
	d := NewDocument()

	if d.State.Value == nil {
		t.Fatal("d.State.Value == nil")
	}
	if d.Events == nil {
		t.Fatal("d.Events == nil")
	}
	if d.EventsByParent == nil {
		t.Fatal("d.EventsByParent == nil")
	}
	if d.Quorums == nil {
		t.Fatal("d.Quorums == nil")
	}
	if d.QuorumsByEvent == nil {
		t.Fatal("d.QuorumsByEvent == nil")
	}
}
