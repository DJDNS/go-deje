package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"os"

	"github.com/DJDNS/go-deje/document"
)

var usage_string = `
USAGE:

    # Convert snapshot to flat-history document cache
    djconvert up original.json deje.json
    djconvert up original.json deje.json --pretty

    # Export event in document cache to snapshot
    djconvert down deje.json snapshot.json 89efc6
    djconvert down deje.json snapshot.json 89efc6 --pretty
`

func usage() {
	fmt.Fprint(os.Stderr, usage_string)
	os.Exit(1)
}

func main() {
	log.SetFlags(0)
	log.SetPrefix("djconvert: ")
	args := flag.Args()

	if len(args) < 1 {
		log.Println("Not enough args, need at least a subcommand 'up' or 'down'")
		usage()
	}

	doc := document.NewDocument()
	input_filename := "input.json"
	output_filename := "output.json"

	input, err := os.Open(input_filename)
	if err != nil {
		log.Fatal(err)
	}

	output, err := os.Create(output_filename)
	if err != nil {
		log.Fatal(err)
	}

	var data interface{}
	if err = json.NewDecoder(input).Decode(&data); err != nil {
		log.Fatal(err)
	}

	event := doc.NewEvent("SET")
	event.Arguments["path"] = []interface{}{}
	event.Arguments["value"] = data
	event.Register()

	if err := doc.Serialize(output); err != nil {
		log.Fatal(err)
	}
	log.Printf("Successfully wrote %s\n", output_filename)
}
