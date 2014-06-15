package document

import (
	"reflect"
	"testing"
)

// At some point, will probably want to find a way to make these
// available for all packages.

func compare(t *testing.T, expected, got interface{}) {
	if !reflect.DeepEqual(expected, got) {
		t.Fatalf("Expected %#v, got %#v", expected, got)
	}
}
func comparem(t *testing.T, expected, got interface{}, msg string) {
	if !reflect.DeepEqual(expected, got) {
		t.Fatalf("%s\nExp %#v\n\nGot %#v", msg, expected, got)
	}
}
