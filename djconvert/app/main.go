package app

import (
	"io"
	"log"
	"os"
	"sort"
	"strings"

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

	var pretty bool = args["--pretty"].(bool)

	if args["up"] == true {
		input_filename := args["<source>"].(string)
		output_filename := args["<target>"].(string)

		input, output, err := get_filehandles(input_filename, output_filename)
		if err != nil {
			logger.Fatal(err)
		}
		if err = up(input, output, pretty); err != nil {
			logger.Fatal(err)
		}
		logger.Printf("Successfully wrote %s\n", output_filename)
	} else if args["down"] == true {
		input_filename := args["<source>"].(string)
		output_filename := args["<target>"].(string)
		hash_prefix := args["<event-hash>"].(string)

		input, output, err := get_filehandles(input_filename, output_filename)
		if err != nil {
			logger.Fatal(err)
		}
		if err, doc := down(input, output, hash_prefix, pretty); err != nil {
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
	}
}

func getFilehandles(input_fn, output_fn) (io.Reader, JsonWriter, error) {
	input, err := os.Open(input_fn)
	if err != nil {
		return nil, nil, err
	}

	output, err := os.Create(output_fn)
	if err != nil {
		return nil, nil, err
	}

	return input, NewJsonWriter(output), nil
}
