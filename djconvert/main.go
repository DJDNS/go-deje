package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"log"
	"os"

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

	if err := doc.Serialize(output); err != nil {
		return err
	}
	return nil
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
	} else {
		log.Printf("Unknown subcommand '%s'\n", subcommand)
		usage()
	}
}
