package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"os"
	"sort"
	"strings"

	"github.com/DJDNS/go-deje/document"
	"github.com/docopt/docopt-go"
)

var version = "djconvert 0.0.13"
var usage_string = `djconvert - Converts files to and from DocCache format.

Usage:
    djconvert up <source> <target> [--pretty]
    djconvert down <source> <target> <event-hash> [--pretty]
    djconvert -h | --help
    djconvert --version

Options:
    --pretty      Pretty-print JSON in output file.
    -h --help     Show this message.
    --version     Show version info.
`

func get_filehandles(input_fn, output_fn string) (io.Reader, io.Writer, error) {
	input, err := os.Open(input_fn)
	if err != nil {
		return nil, nil, err
	}

	output, err := os.Create(output_fn)
	if err != nil {
		return nil, nil, err
	}

	return input, output, nil
}

func write_json(data interface{}, output io.Writer, pretty bool) error {
	if pretty {
		// Go JSON API is a bit clumsy :(
		buf, err := json.MarshalIndent(data, "", "    ")
		if err != nil {
			return err
		}
		n, err := fmt.Fprint(output, string(buf))
		if err != nil {
			return err
		}
		if n < len(buf) {
			return errors.New("Didn't write all bytes to file")
		}
		return nil
	} else {
		return json.NewEncoder(output).Encode(&data)
	}
}

func up(input io.Reader, output io.Writer, pretty bool) error {
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

	return write_json(doc, output, pretty)
}

func down(input io.Reader, output io.Writer, hash_prefix string, pretty bool) (error, document.Document) {
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
			return errors.New("No such hash '" + hash_prefix + "'"), doc
		} else if len(found) > 1 {
			return errors.New("Hash prefix '" + hash_prefix + "' is ambiguous"), doc
		} else {
			event = found[0]
		}
	}

	if err := event.Goto(); err != nil {
		return err, doc
	}
	return write_json(doc.State.Export(), output, pretty), doc
}

func main() {
	log.SetFlags(0)
	log.SetPrefix("djconvert: ")

	args, err := docopt.Parse(usage_string, nil, true, version, false, true)
	if err != nil {
		log.Fatal(err)
	}

	var pretty bool = args["--pretty"].(bool)

	if args["up"] == true {
		input_filename := args["<source>"].(string)
		output_filename := args["<target>"].(string)

		input, output, err := get_filehandles(input_filename, output_filename)
		if err != nil {
			log.Fatal(err)
		}
		if err = up(input, output, pretty); err != nil {
			log.Fatal(err)
		}
		log.Printf("Successfully wrote %s\n", output_filename)
	} else if args["down"] == true {
		input_filename := args["<source>"].(string)
		output_filename := args["<target>"].(string)
		hash_prefix := args["<hash-prefix>"].(string)

		input, output, err := get_filehandles(input_filename, output_filename)
		if err != nil {
			log.Fatal(err)
		}
		if err, doc := down(input, output, hash_prefix, pretty); err != nil {
			log.Print(err)

			keys := make([]string, len(doc.Events))
			var i int
			for key := range doc.Events {
				keys[i] = key
				i++
			}
			sort.Strings(keys)

			log.Fatalf("Available hashes (%d):\n%s",
				len(doc.Events),
				strings.Join(keys, "\t\n"),
			)
		}
	}
}
