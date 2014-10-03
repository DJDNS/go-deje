package app

import (
	"io"
	"log"
	"os"

	"github.com/docopt/docopt.go"
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

func Main(argv []string, exit bool, log_writer io.Writer) {
	logger := log.New(log_writer, "djconvert: ", 0)

	args, err := docopt.Parse(usage_string, nil, true, version, false, true)
	if err != nil {
		logger.Fatal(err)
	}

	input_filename := args["<source>"].(string)
	output_filename := args["<target>"].(string)
	pretty := args["--pretty"].(bool)

	input, output, err := getFilehandles(input_filename, output_filename, pretty)
	if err != nil {
		logger.Fatal(err)
	}

	var command func() error

	if args["up"] == true {
		command = func() error {
			return DoCommandUp(input, output)
		}
	} else if args["down"] == true {
		hash_prefix := args["<event-hash>"].(string)
		command = func() error {
			return DoCommandDown(input, output, hash_prefix)
		}
		/*
			if err := DoCommandDown(input, output, hash_prefix); err != nil {
				logger.Print(err)

				keys := make([]string, len(doc.Events))
				var i int
				for key := range doc.Events {
					keys[i] = key
					i++
				}
				sort.Strings(keys)

				logger.Fatalf("Available hashes (%d):\n%s",
					len(doc.Events),
					strings.Join(keys, "\t\n"),
				)
			}
		*/
	}
	if err := command(); err != nil {
		logger.Fatal(err)
	}
	logger.Printf("Successfully wrote %s\n", output_filename)
}

func getFilehandles(input_fn, output_fn string, pretty bool) (io.Reader, JsonWriter, error) {
	input, err := os.Open(input_fn)
	if err != nil {
		return nil, nil, err
	}

	output, err := os.Create(output_fn)
	if err != nil {
		return nil, nil, err
	}

	return input, NewJsonWriter(output, pretty), nil
}
