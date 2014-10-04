package app

import (
	"bytes"
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
