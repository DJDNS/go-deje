package model

import "testing"

func TestNewDocumentFile(t *testing.T) {
	df := NewDocumentFile()

	if df.Events == nil {
		t.Fatal("df.Events not initialized")
	}

	if df.Quorums == nil {
		t.Fatal("df.Quorums not initialized")
	}
}
