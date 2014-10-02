package app

import (
	"encoding/json"
	"io"

	"github.com/DJDNS/go-deje/document"
)

func DoCommandUp(input io.Reader, output JsonWriter) error {
	doc := document.NewDocument()

	var data interface{}
	if err := json.NewDecoder(input).Decode(&data); err != nil {
		return err
	}

	event := doc.NewEvent("SET")
	event.Arguments["path"] = []interface{}{}
	event.Arguments["value"] = data
	event.Register()

	doc.Timestamps = append(doc.Timestamps, event.Hash())

	return output.Write(doc)
}
