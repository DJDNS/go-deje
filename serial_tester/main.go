package main

import (
	"encoding/json"
	"errors"
	"io"
	"log"
	"os"

	"github.com/DJDNS/go-deje/document"
	"github.com/DJDNS/go-deje/util"
)

func deserializeDocument(r io.Reader) (interface{}, error) {
	doc := document.NewDocument()
	err := doc.Deserialize(r)
	return doc, err
}
func deserializeEvent(r io.Reader) (interface{}, error) {
	event := document.Event{}
	decoder := json.NewDecoder(r)
	err := decoder.Decode(&event)
	return event, err
}
func deserializeInput(object_type string, r io.Reader) (interface{}, error) {
	switch object_type {
	case "document":
		return deserializeDocument(r)
	case "event":
		return deserializeEvent(r)
	default:
		return nil, errors.New("No such object type: " + object_type)
	}
}

func formatPretty4(object interface{}, w io.Writer) error {
	b, err := json.MarshalIndent(object, "", "    ")
	if err != nil {
		return err
	}
	_, err = w.Write(b)
	return err
}
func formatCompact(object interface{}, w io.Writer) error {
	encoder := json.NewEncoder(w)
	return encoder.Encode(object)
}
func formatHash(object interface{}, w io.Writer) error {
	hash, err := util.HashObject(object)
	if err != nil {
		return err
	}
	_, err = w.Write([]byte(hash))
	return err
}
func serializeOutput(format string, object interface{}, w io.Writer) error {
	switch format {
	case "pretty4":
		return formatPretty4(object, w)
	case "compact":
		return formatCompact(object, w)
	case "hash":
		return formatHash(object, w)
	default:
		return errors.New("No such format: " + format)
	}
}

func main() {
	args := os.Args
	if len(args) < 3 {
		log.Fatal("Insufficient arguments")
	}

	object_type := args[1]
	format := args[2]

	object, err := deserializeInput(object_type, os.Stdin)
	if err != nil {
		log.Fatal(err)
	}

	err = serializeOutput(format, object, os.Stdout)
	if err != nil {
		log.Fatal(err)
	}
}
