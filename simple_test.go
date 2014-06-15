package deje

import "testing"

func TestSimpleClient_NewSimpleClient(t *testing.T) {
	topic := "http://example.com/deje/some-doc"
	sc := NewSimpleClient(topic)
	if sc.client.Doc.Topic != topic {
		t.Fatal("Did not create encapsulated Client correctly")
	}
}
