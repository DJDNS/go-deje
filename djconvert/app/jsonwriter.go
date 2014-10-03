package app

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
)

type JsonWriter interface {
	Write(data interface{}) error
}

func NewJsonWriter(w io.Writer, pretty bool) JsonWriter {
	if pretty {
		return PrettyJsonWriter{w}
	} else {
		return CompactJsonWriter{w}
	}
}

type PrettyJsonWriter struct {
	w io.Writer
}

func (jw PrettyJsonWriter) Write(data interface{}) error {
	// Go JSON API is a bit clumsy :(
	buf, err := json.MarshalIndent(data, "", "    ")
	if err != nil {
		return err
	}
	n, err := fmt.Fprint(jw.w, string(buf))
	if err != nil {
		return err
	}
	if n < len(buf) {
		return errors.New("Didn't write all bytes to file")
	}
	return nil
}

type CompactJsonWriter struct {
	w io.Writer
}

func (jw CompactJsonWriter) Write(data interface{}) error {
	return json.NewEncoder(jw.w).Encode(&data)
}
