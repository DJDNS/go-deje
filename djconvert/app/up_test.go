package app

import (
	"bytes"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestDoCommandUp(t *testing.T) {
	buf := new(bytes.Buffer)
	input_data := `{
		"foo": "bar"
	}`
	reader := strings.NewReader(input_data)
	writer := NewJsonWriter(buf, true)

	expected_output := `{
    "events": {
        "41a772b775b1c4afbbfb42b7a91b3031a712ab42": {
            "parent": "",
            "handler": "SET",
            "args": {
                "path": [],
                "value": {
                    "foo": "bar"
                }
            }
        }
    },
    "timestamps": [
        "41a772b775b1c4afbbfb42b7a91b3031a712ab42"
    ]
}`
	assert.NoError(t, DoCommandUp(reader, writer))
	assert.Equal(t, expected_output, buf.String())
}

func TestDoCommandUp_BadInput(t *testing.T) {
	buf := new(bytes.Buffer)
	input_data := `{ lol {`
	reader := strings.NewReader(input_data)
	writer := NewJsonWriter(buf, true)

	assert.Error(t, DoCommandUp(reader, writer))
	assert.Equal(t, "", buf.String())
}
