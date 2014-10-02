package app

import (
	"errors"
	"io"
	"strings"

	"github.com/DJDNS/go-deje/document"
)

func DoCommandDown(input io.Reader, output JsonWriter, hash_prefix string) error {
	doc := document.NewDocument()
	doc.Deserialize(input)

	event, ok := doc.Events[hash_prefix]
	if !ok {
		found := make([]*document.Event, 0)
		for key := range doc.Events {
			if strings.HasPrefix(key, hash_prefix) {
				found = append(found, doc.Events[key])
			}
		}
		if len(found) == 0 {
			return errors.New("No such hash '" + hash_prefix + "'")
		} else if len(found) > 1 {
			return errors.New("Hash prefix '" + hash_prefix + "' is ambiguous")
		} else {
			event = found[0]
		}
	}

	if err := event.Goto(); err != nil {
		return err
	}
	return output.Write(doc.State.Export())
}
