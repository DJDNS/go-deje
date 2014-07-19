package main

import (
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"io"
	"log"
	"os"
	"sort"
	"strings"

	"github.com/DJDNS/go-deje/document"
)

var usage_string = `
USAGE:

    # Convert snapshot to flat-history document cache
    djconvert up original.json deje.json
    djconvert --pretty up original.json deje.json

    # Export event in document cache to snapshot
    djconvert down deje.json snapshot.json 89efc6
    djconvert --pretty down deje.json snapshot.json 89efc6
`

func usage() {
	fmt.Fprint(os.Stderr, usage_string)
	os.Exit(1)
}

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

func write_json(data interface{}, output io.Writer) error {
	if *pretty {
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

func up(input io.Reader, output io.Writer) error {
	doc := document.NewDocument()

	var data interface{}
	if err := json.NewDecoder(input).Decode(&data); err != nil {
		return err
	}

	event := doc.NewEvent("SET")
	event.Arguments["path"] = []interface{}{}
	event.Arguments["value"] = data
	event.Register()

	return write_json(doc, output)
}

func down(input io.Reader, output io.Writer, hash_prefix string) (error, document.Document) {
	doc := document.NewDocument()
	doc.Deserialize(input)

	event, ok := doc.Events[hash_prefix]
	if !ok {
		return errors.New("No such hash '" + hash_prefix + "'"), doc
	}

	if err := event.Goto(); err != nil {
		return err, doc
	}
	return write_json(doc.State.Export(), output), doc
}

var pretty = flag.Bool("pretty", false, "Pretty-print the output, for human readability")

func main() {
	log.SetFlags(0)
	log.SetPrefix("djconvert: ")
	flag.Parse()
	args := flag.Args()

	if len(args) < 1 {
		log.Println("Not enough args, need at least a subcommand 'up' or 'down'")
		usage()
	}

	subcommand := args[0]
	if subcommand == "up" {
		if len(args) < 3 {
			log.Println("Subcommand 'up' takes 2 additional args")
			usage()
		}
		input_filename, output_filename := args[1], args[2]
		input, output, err := get_filehandles(input_filename, output_filename)
		if err != nil {
			log.Fatal(err)
		}
		if err = up(input, output); err != nil {
			log.Fatal(err)
		}
		log.Printf("Successfully wrote %s\n", output_filename)
	} else if subcommand == "down" {
		if len(args) < 4 {
			log.Println("Subcommand 'down' takes 3 additional args")
			usage()
		}
		input_filename, output_filename, hash_prefix := args[1], args[2], args[3]
		input, output, err := get_filehandles(input_filename, output_filename)
		if err != nil {
			log.Fatal(err)
		}
		if err, doc := down(input, output, hash_prefix); err != nil {
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
	} else {
		log.Printf("Unknown subcommand '%s'\n", subcommand)
		usage()
	}
}
