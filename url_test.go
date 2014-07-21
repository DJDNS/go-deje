package deje

import "testing"
import "github.com/stretchr/testify/assert"

type UrlTest struct {
	Input  string
	Router string
	Topic  string
}

func (ut UrlTest) Test(t *testing.T) {
	router, topic, err := GetRouterAndTopic(ut.Input)
	if err != nil {
		router = "<error>: " + err.Error()
		topic = ""
	}
	assert.Equal(t, ut.Router, router, "Router output is what was expected (%s)", ut.Input)
	assert.Equal(t, ut.Topic, topic, "Topic output is what was expected (%s)", ut.Input)
}

func TestGetRouterAndTopic(t *testing.T) {
	tests := []UrlTest{
		UrlTest{
			Input:  "deje://foo/bar",
			Router: "ws://foo/ws",
			Topic:  "deje://foo/bar",
		},
		UrlTest{
			Input:  "deje://foo:8080",
			Router: "ws://foo:8080/ws",
			Topic:  "deje://foo:8080/",
		},
		UrlTest{
			Input:  "foo.bar.baz",
			Router: "<error>: URL does not start with 'deje://': 'foo.bar.baz'",
			Topic:  "",
		},
		UrlTest{
			Input:  "deje://foo.bar.baz:8080",
			Router: "ws://foo.bar.baz:8080/ws",
			Topic:  "deje://foo.bar.baz:8080/",
		},
		UrlTest{
			Input:  "//foo.bar.baz:8080",
			Router: "<error>: URL does not start with 'deje://': '//foo.bar.baz:8080'",
			Topic:  "",
		},
		UrlTest{
			Input:  "deje://foo.bar.baz:8080",
			Router: "ws://foo.bar.baz:8080/ws",
			Topic:  "deje://foo.bar.baz:8080/",
		},
		UrlTest{
			Input:  "deje://%",
			Router: "<error>: parse deje://%: hexadecimal escape in host",
			Topic:  "",
		},
	}
	for _, ut := range tests {
		ut.Test(t)
	}
}
