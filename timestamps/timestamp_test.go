package timestamps

import (
	"github.com/DJDNS/go-deje/document"
	"reflect"
	"testing"
)

func TestDTS_GetTimestamps(t *testing.T) {
	dts := DummyTimestampService{}
	stamps, err := dts.GetTimestamps("Interstella")
	if len(stamps) != 0 {
		t.Fatalf(
			"Expected empty timestamp array, has %d elements",
			len(stamps),
		)
	}
	if err != nil {
		t.Fatal(err)
	}
}

func TestNewSTS(t *testing.T) {
	doc := document.NewDocument()
	doc.Topic = "furbies"
	sts := NewSortingTimestampService(doc)
	if !reflect.DeepEqual(sts.Doc, doc) {
		t.Fatalf("%#v != %#v", sts.Doc, doc)
	}
}
func TestSTS_GetTimestamps(t *testing.T) {
	doc := document.NewDocument()
	doc.Topic = "furbies"
	sts := NewSortingTimestampService(doc)

	// The given values are event hashes. The quorum hashes will be
	// different than these string literals.
	for _, evhash := range []string{"123", "456", "789"} {
		q := doc.NewQuorum(evhash)
		q.Register()
	}

	timestamps, err := sts.GetTimestamps(doc.Topic)
	if err != nil {
		t.Fatal(err)
	}
	expected_timestamps := []string{
		"5d3b9fa37c8145112882e77b1aa5db9477dab734",
		"e211156f9d2c736a6d1718246216f97974ca9585",
		"fce5d4ea3a4e2c130657bf97b286b16f54da6850",
	}
	if !reflect.DeepEqual(timestamps, expected_timestamps) {
		t.Fatalf("Expected %#v, got %#v", expected_timestamps, timestamps)
	}
}
