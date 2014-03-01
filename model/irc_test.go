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

func TestIRCLocation_GetVariables(t *testing.T) {
	location := IRCLocation{"example.com", 9999, "thechannel"}
	vars := location.GetVariables()
	expected := "host=example.com&port=9999&channel=thechannel"

	if vars != expected {
		t.Fatalf("Expected %s, got %s", expected, vars)
	}
}

type loc_parse_test struct {
	URL    string
	Result IRCLocation
}

var parse_tests = []loc_parse_test{
	loc_parse_test{
		"irc://hello.world/",
		IRCLocation{"hello.world", 6667, ""},
	},
	loc_parse_test{
		"irc://hello.world:32/#blah",
		IRCLocation{"hello.world", 32, "blah"},
	},
}

func TestIRCLocation_ParseFrom(t *testing.T) {
	location := IRCLocation{}

	err := location.ParseFrom("%")
	if err == nil {
		t.Fatal("ParseFrom should reject malformed URLs")
	}

	err = location.ParseFrom("irc://host:bad_port/")
	if err == nil {
		t.Fatal("ParseFrom should reject non-int port")
	}

	for _, test := range parse_tests {
		err = location.ParseFrom(test.URL)
		if err != nil {
			t.Fatal(err)
		}
		if location != test.Result {
			t.Fatalf("Expected %v, got %v", test.Result, location)
		}
	}

}
