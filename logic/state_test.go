package logic

import "testing"

func TestGetProperty(t *testing.T) {
	ds := NewDocumentState()
	ds.Content["hello"] = "world"

	var mystr string
	err := ds.GetProperty("hello", &mystr)
	if err != nil {
		t.Fatalf("GetProperty failed: %v", err)
	}

	if mystr != "world" {
		t.Fatal("GetProperty did not retrieve value")
	}
}

func TestGetProperty_Missing(t *testing.T) {
	ds := NewDocumentState()

	var dummy []int
	err := ds.GetProperty("stuff", dummy)

	if err == nil {
		t.Fatal("GetProperty should have failed")
	}
}
