package app

import (
	"bytes"
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestJsonWriter(t *testing.T) {
	kabuki := []interface{}{
		8, "kabuki", false,
		map[string]interface{}{"bobby": "singer"},
	}
	tests := []struct {
		Pretty         bool
		Data           interface{}
		ExpectedOutput string
	}{
		{
			true,
			kabuki,
			`[
    8,
    "kabuki",
    false,
    {
        "bobby": "singer"
    }
]` + "\n",
		},
		{
			false,
			kabuki,
			"[8,\"kabuki\",false,{\"bobby\":\"singer\"}]\n",
		},
	}
	for _, test := range tests {
		buf := new(bytes.Buffer)
		jw := NewJsonWriter(buf, test.Pretty)
		err := jw.Write(test.Data)

		assert.NoError(t, err)
		assert.Equal(t, test.ExpectedOutput, buf.String())
	}

}

func TestPrettyJsonWriter_Fail_Serialization(t *testing.T) {
	buf := new(bytes.Buffer)
	jw := NewJsonWriter(buf, true)
	err := jw.Write(make(chan int))
	if assert.Error(t, err) {
		assert.Equal(t, "json: unsupported type: chan int", err.Error())
	}
}

type FakeWriter int

func (fw FakeWriter) Write([]byte) (int, error) {
	num_bytes := int(fw)
	if num_bytes == 0 {
		return 0, errors.New("Fails immediately")
	} else {
		return num_bytes, nil
	}
}

func TestPrettyJsonWriter_Fail_Fprintf(t *testing.T) {
	tests := []struct {
		Writer   FakeWriter
		ErrorMsg string
	}{
		{
			FakeWriter(0),
			"Fails immediately",
		},
		{
			FakeWriter(5),
			"Didn't write all bytes to file",
		},
	}
	for _, test := range tests {
		jw := NewJsonWriter(test.Writer, true)
		err := jw.Write("Big long JSON string")
		if assert.Error(t, err) {
			assert.Equal(t, test.ErrorMsg, err.Error())
		}
	}
}
