package model

import "testing"

func TestIRCLocation_GetURL(t *testing.T) {
	location := IRCLocation{"example.com", 9999, "thechannel"}
	url := location.GetURL()
	expected := "irc://example.com:9999/#thechannel"

	if url != expected {
		t.Fatalf("Expected %s, got %s", expected, url)
	}
}
