package app

import (
	"fmt"
	"io"
	"sort"
	"strings"

	"github.com/DJDNS/go-deje/document"
)

type HashPrefixError struct {
	Problem string
	Prefix  string
	Doc     *document.Document
}

// Return a sorted list of event hashes for printing.
func (hpe HashPrefixError) GetEventHashesAvailable() []string {
	keys := make([]string, len(hpe.Doc.Events))
	var i int
	for key := range hpe.Doc.Events {
		keys[i] = key
		i++
	}
	sort.Strings(keys)

	return keys
}

func (hpe HashPrefixError) Error() string {
	keys := hpe.GetEventHashesAvailable()

	return fmt.Sprintf("%s: %s\n\nAvailable hashes (%d):\n%s",
		hpe.Problem, hpe.Prefix,
		len(keys), strings.Join(keys, "\t\n"),
	)
}

func GetEventByPrefix(doc *document.Document, hash_prefix string) (*document.Event, error) {
	event, ok := doc.Events[hash_prefix]
	if ok {
		return event, nil
	}

	// Not an exact match - do prefix search
	found := make([]*document.Event, 0)
	for key := range doc.Events {
		if strings.HasPrefix(key, hash_prefix) {
			found = append(found, doc.Events[key])
		}
	}

	// Determine return value by number of results
	if len(found) == 0 {
		return nil, HashPrefixError{"No such event", hash_prefix, doc}
	} else if len(found) > 1 {
		return nil, HashPrefixError{"Ambiguous hash prefix", hash_prefix, doc}
	} else {
		return found[0], nil
	}
}

func DoCommandDown(input io.Reader, output JsonWriter, hash_prefix string) error {
	doc := document.NewDocument()
	doc.Deserialize(input)

	event, err := GetEventByPrefix(&doc, hash_prefix)
	if err != nil {
		return err
	}

	if err = event.Goto(); err != nil {
		return err
	}
	return output.Write(doc.State.Export())
}
