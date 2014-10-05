package app

import (
	"io"
	"os"

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
var do_help = true
var options_first = false

func Main(argv []string, exit bool) error {
	do_help = exit
	args, err := docopt.Parse(usage_string, argv, do_help, version, options_first, exit)
	if err != nil {
		return err
	}

	input_filename := args["<source>"].(string)
	output_filename := args["<target>"].(string)
	pretty := args["--pretty"].(bool)

	input, output, err := getFilehandles(input_filename, output_filename, pretty)
	if err != nil {
		return err
	}

	if args["up"] == true {
		return DoCommandUp(input, output)
	} else {
		hash_prefix := args["<event-hash>"].(string)
		return DoCommandDown(input, output, hash_prefix)
	}
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
