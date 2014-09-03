package timestamps

import (
	"reflect"
	"testing"

	"github.com/DJDNS/go-deje/document"
	"github.com/stretchr/testify/assert"
)

func TestDTS_GetTimestamps(t *testing.T) {
	dts := DummyTimestampService{}
	stamps, err := dts.GetTimestamps()
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
	sts := NewSortingTimestampService(doc)
	if !reflect.DeepEqual(sts.Doc, doc) {
		t.Fatalf("%#v != %#v", sts.Doc, doc)
	}
}
func TestSTS_GetTimestamps(t *testing.T) {
	doc := document.NewDocument()
	sts := NewSortingTimestampService(doc)

	// The given values are event "types".
	// Hashes will be different than these string literals.
	for _, evhash := range []string{"123", "456", "789"} {
		q := doc.NewEvent(evhash)
		q.Register()
	}

	timestamps, err := sts.GetTimestamps()
	if err != nil {
		t.Fatal(err)
	}
	expected_timestamps := []string{
		"2303adf72049c8f0d2dd3c38d47775f9e0b0458d",
		"5d0d0d82f38428c33802403af6fdf27e82fcd4bc",
		"f35ae012679b73922225d21834bf962f2c8f1145",
	}
	if !reflect.DeepEqual(timestamps, expected_timestamps) {
		t.Fatalf("Expected %#v, got %#v", expected_timestamps, timestamps)
	}
}

func TestNewPTS(t *testing.T) {
	doc := document.NewDocument()
	sts := NewPeerTimestampService(&doc)
	assert.Equal(t, &doc, sts.Doc)
}
func TestPTS_GetTimestamps(t *testing.T) {
	doc := document.NewDocument()
	sts := NewPeerTimestampService(&doc)

	doc.Timestamps = []string{"123", "456", "789"}

	timestamps, err := sts.GetTimestamps()
	if err != nil {
		t.Fatal(err)
	}
	assert.Equal(t, doc.Timestamps, timestamps)
}
