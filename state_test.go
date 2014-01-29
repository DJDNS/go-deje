package deje

import (
	"reflect"
	"testing"
)

func TestDSGetChannel(t *testing.T) {
	ds := NewDocumentState()

	channel := make(JSONObject)
	channel["host"] = "some string"
	channel["port"] = 9001
	channel["channel"] = "go-nuts"

	ds.Content["channel"] = channel

	loc, err := ds.GetChannel()
	if err != nil {
		t.Fatal(err)
	}

	expected := IRCLocation{
		Host:    "some string",
		Port:    9001,
		Channel: "go-nuts",
	}
	if *loc != expected {
		t.Fatalf("Expected %v, got %v", expected, loc)
	}
}

func TestDSGetChannelBadData(t *testing.T) {
	ds := NewDocumentState()

	_, err := ds.GetChannel()
	if err == nil {
		t.Fatal("GetChannel should have failed, but didn't")
	}

	ds.Content["channel"] = 4
	_, err = ds.GetChannel()
	if err == nil {
		t.Fatal("GetChannel should have failed, but didn't")
	}

	channel := make(JSONObject)
	channel["port"] = "string port"
	ds.Content["channel"] = channel
	_, err = ds.GetChannel()
	if err == nil {
		t.Fatal("GetChannel should have failed, but didn't")
	}
}

func TestDSGetDownloads(t *testing.T) {
	ds := NewDocumentState()

	urls := []interface{}{"a", "b", "c"}
	ds.Content["urls"] = urls

	got, err := ds.GetURLs()
	if err != nil {
		t.Fatal(err)
	}

	expected := DownloadURLs{"a", "b", "c"}
	if !reflect.DeepEqual(*got, expected) {
		t.Fatalf("Expected %v, got %v", expected, got)
	}
}

func TestDSGetDownloadsBadData(t *testing.T) {
	ds := NewDocumentState()

	_, err := ds.GetURLs()
	if err == nil {
		t.Fatal("GetURLs should have failed, but didn't")
	}

	ds.Content["urls"] = make(JSONObject)
	_, err = ds.GetURLs()
	if err == nil {
		t.Fatal("GetURLs should have failed, but didn't")
	}

	ds.Content["urls"] = []interface{}{1, 2, 3}
	_, err = ds.GetURLs()
	if err == nil {
		t.Fatal("GetURLs should have failed, but didn't")
	}
}
