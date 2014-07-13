package main

import (
	"encoding/json"
	"log"
	"os"

	"github.com/DJDNS/go-deje/document"
)

func main() {
	log.SetFlags(0)
	log.SetPrefix("djconvert: ")

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
	log.Printf("Successfully wrote %s", output_filename)
}
