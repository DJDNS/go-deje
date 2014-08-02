package timestamps

import (
	"errors"
	"testing"

	"github.com/DJDNS/go-deje/document"
	"github.com/stretchr/testify/assert"
)

func TestTimestampTracker_StartIteration(t *testing.T) {
	doc := document.NewDocument()
	service := NewPeerTimestampService(&doc)
	tracker := NewTimestampTracker(&doc, service)

	tracker.timestamps = []string{"1", "2", "3"}
	tracker.tip = "marshmallow"

	assert.NoError(t, tracker.StartIteration())
	assert.Equal(t, doc.Timestamps, tracker.timestamps)
	assert.Equal(t, "", tracker.tip)
}

type failingTimestampService string

func (fts failingTimestampService) GetTimestamps(topic string) ([]string, error) {
	return nil, errors.New(string(fts))
}

func TestTimestampTracker_StartIteration_ServiceFailure(t *testing.T) {
	doc := document.NewDocument()
	service := failingTimestampService("Failure message")
	tracker := NewTimestampTracker(&doc, service)

	if err := tracker.StartIteration(); assert.Error(t, err) {
		assert.Equal(t, "Failure message", err.Error())
	}
}
