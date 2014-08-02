package timestamps

import (
	"errors"
	"log"

	"github.com/DJDNS/go-deje/document"
)

type TimestampTracker struct {
	Doc     *document.Document
	Service TimestampService

	// Current iteration range
	timestamps []string
	tip        string
}

func NewTimestampTracker(doc *document.Document, service TimestampService) TimestampTracker {
	return TimestampTracker{
		Doc:     doc,
		Service: service,
	}
}

// Iterate until tip is found. Returns tip
func (tt *TimestampTracker) GoToLatest(logger *log.Logger) string {
	if err := tt.StartIteration(); err != nil {
		if logger != nil {
			logger.Println(err)
		}
		return ""
	}
	for p := range tt.timestamps {
		err := tt.DoIteration(p)
		if err != nil && logger != nil {
			logger.Printf("Error on iteration %d (current tip: '%s'):\n", p, tt.tip)
			logger.Println(err)
		}
	}
	return tt.tip
}

// Set up to find tip event - reset finder state.
func (tt *TimestampTracker) StartIteration() error {
	timestamps, err := tt.Service.GetTimestamps(tt.Doc.Topic)
	if err != nil {
		return err
	}

	tt.timestamps = timestamps
	tt.tip = ""
	return nil
}

// Single iteration of finder. If it succeeds, updates the tip.
func (tt *TimestampTracker) DoIteration(position int) error {
	if position < 0 || position > len(tt.timestamps) {
		return errors.New("Bad position")
	}

	ts := tt.timestamps[position]
	event, ok := tt.Doc.Events[ts]
	if !ok {
		return errors.New("No such event")
	}

	if !tt.CompatibleWithTip(event) {
		return errors.New("Event is not compatible with and ahead of tip")
	}

	if err := event.Goto(); err != nil {
		return err
	}

	// Timestamp has passed through the gauntlet successfully
	tt.tip = ts
	return nil
}

// TODO: Implement correctly
func (tt *TimestampTracker) CompatibleWithTip(event *document.Event) bool {
	return false
}
