package app

import (
	"bytes"
	"strings"
	"testing"

	"github.com/DJDNS/go-deje/document"
	"github.com/stretchr/testify/assert"
)

func TestGetEventByPrefix(t *testing.T) {
	doc_no_events := document.NewDocument()
	doc_with_events := document.NewDocument()

	// Contents don't really matter, just need variety of hashes
	evA := doc_with_events.NewEvent("A")
	evA.Register()
	evB := doc_with_events.NewEvent("B")
	evB.Register()
	evC := doc_with_events.NewEvent("C")
	evC.Register()
	evD := doc_with_events.NewEvent("Common prefix with evC")
	evD.Register()

	sorted_hashes := evA.Hash() + "\n\t" +
		evC.Hash() + "\n\t" +
		evD.Hash() + "\n\t" +
		evB.Hash()

	tests := []struct {
		Doc           *document.Document
		HashPrefix    string
		ExpectedEvent *document.Event
		ExpectedError string
	}{
		// Empty hash prefix
		{
			&doc_no_events, "",
			nil,
			"No such event: ''\n\nAvailable hashes (0):\n",
		},
		// No such event
		{
			&doc_no_events, "foo",
			nil,
			"No such event: 'foo'\n\nAvailable hashes (0):\n",
		},
		// No such event (where events are available)
		{
			&doc_with_events, "foo",
			nil,
			"No such event: 'foo'\n\nAvailable hashes (4):\n\t" + sorted_hashes,
		},
		// Exact match
		{
			&doc_with_events, evC.Hash(),
			&evC,
			"",
		},
		// Valid prefix
		{
			&doc_with_events, "e3a",
			&evB,
			"",
		},
		// Colliding prefix
		{
			&doc_with_events, "e",
			nil,
			"Ambiguous hash prefix: 'e'\n\nAvailable hashes (4):\n\t" + sorted_hashes,
		},
	}
	for _, test := range tests {
		event, err := GetEventByPrefix(test.Doc, test.HashPrefix)
		assert.Equal(t, test.ExpectedEvent, event)
		if test.ExpectedError == "" {
			assert.NoError(t, err)
		} else {
			if assert.Error(t, err) {
				assert.Equal(t, test.ExpectedError, err.Error())
			}
		}
	}
}

func TestDoCommandDown(t *testing.T) {
	tests := []struct {
		Input          string
		HashPrefix     string
		ExpectedOutput string
		ExpectedError  string
	}{
		// Bad input
		{
			"{{",
			"",
			"",
			"invalid character '{' looking for beginning of object key string",
		},
		// Error in GetEventByPrefix
		{
			"{}",
			"Blasternaut",
			"",
			"No such event: 'Blasternaut'\n\nAvailable hashes (0):\n",
		},
		// Bad event (type "SAT", not "SET")
		{
			`{"events":{ "":{"handler":"SAT"} }}`,
			"beac",
			"",
			"Custom events are not supported yet",
		},
		// Success
		{
			`{"events":{ "":{
				"handler":"SET", "args":{
					"path":["hello"], "value":"world"}
				}
			}}`,
			"ff",
			`{"hello":"world"}` + "\n",
			"",
		},
	}
	for _, test := range tests {
		buf := new(bytes.Buffer)
		output_writer := NewJsonWriter(buf, false)
		err := DoCommandDown(
			strings.NewReader(test.Input),
			output_writer,
			test.HashPrefix,
		)

		assert.Equal(t, test.ExpectedOutput, buf.String())
		if test.ExpectedError == "" {
			assert.NoError(t, err)
		} else {
			if assert.Error(t, err) {
				assert.Equal(t, test.ExpectedError, err.Error())
			}
		}
	}
}
