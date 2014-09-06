package timestamps

import "github.com/DJDNS/go-deje/document"

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
func (tt *TimestampTracker) FindLatest() (*document.Event, error) {
	timestamps, err := tt.Service.GetTimestamps()
	if err != nil {
		return nil, err
	}
	tt.timestamps = timestamps
	tt.tip = ""

	for _, ts := range tt.timestamps {
		event, ok := tt.Doc.Events[ts]
		if !ok {
			continue
		}

		if !tt.CompatibleWithTip(event) {
			continue
		}

		tt.tip = ts
	}
	return tt.Doc.Events[tt.tip], nil
}

func (tt *TimestampTracker) CompatibleWithTip(event *document.Event) bool {
	if event == nil {
		return false
	}
	hash := event.Hash()
	if _, ok := tt.Doc.Events[hash]; !ok {
		return false
	}

	// Always win against no-tip
	if tt.tip == "" {
		return true
	}

	tip_event, ok := tt.Doc.Events[tt.tip]
	if !ok {
		return false
	}

	ca, err := tip_event.GetCommonAncestor(event)
	if err != nil {
		return false
	}

	return ca == tip_event
}
