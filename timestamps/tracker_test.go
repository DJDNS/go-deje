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

type compatibleTest struct {
	Tip            string
	ComparedEvent  *document.Event
	ExpectedResult bool
	Description    string
}

func TestTimestampTracker_CompatibleWithTip(t *testing.T) {
	doc := document.NewDocument()
	tracker := NewTimestampTracker(&doc, nil)

	evRoot := doc.NewEvent("SET")
	evRoot.Arguments["path"] = []interface{}{"key"}
	evRoot.Arguments["value"] = "value"
	evRoot.Register()

	evChild := doc.NewEvent("SET")
	evChild.Arguments["path"] = []interface{}{"other key"}
	evChild.Arguments["value"] = "other value"
	evChild.SetParent(evRoot)
	evChild.Register()

	// Competes with evChild
	evFork := doc.NewEvent("SET")
	evFork.Arguments["path"] = []interface{}{"fork"}
	evFork.Arguments["value"] = "fork"
	evFork.SetParent(evRoot)
	evFork.Register()

	evBadParentHash := doc.NewEvent("SET")
	evBadParentHash.ParentHash = "foobarbaz"
	evBadParentHash.Register()

	evUnregistered := doc.NewEvent("foo")

	tests := []compatibleTest{
		compatibleTest{
			Tip:            "",
			ComparedEvent:  nil,
			ExpectedResult: false,
			Description:    "A nil event pointer",
		},
		compatibleTest{
			Tip:            "",
			ComparedEvent:  &evUnregistered,
			ExpectedResult: false,
			Description:    "An unregistered event",
		},
		compatibleTest{
			Tip:            "",
			ComparedEvent:  &evRoot,
			ExpectedResult: true,
			Description:    "Any registered event vs no-tip",
		},
		compatibleTest{
			Tip:            evRoot.Hash(),
			ComparedEvent:  &evChild,
			ExpectedResult: true,
			Description:    "Child of root event",
		},
		compatibleTest{
			Tip:            evChild.Hash(),
			ComparedEvent:  &evRoot,
			ExpectedResult: false,
			Description:    "Parent of tip event",
		},
		compatibleTest{
			Tip:            evChild.Hash(),
			ComparedEvent:  &evFork,
			ExpectedResult: false,
			Description:    "Incompatible forks",
		},
		compatibleTest{
			Tip:            "foobar",
			ComparedEvent:  &evRoot,
			ExpectedResult: false,
			Description:    "Some random broken tip value",
		},
		compatibleTest{
			Tip:            evRoot.Hash(),
			ComparedEvent:  &evBadParentHash,
			ExpectedResult: false,
			Description:    "Bad or missing heritage",
		},
	}
	for _, test := range tests {
		tracker.tip = test.Tip
		assert.Equal(t,
			test.ExpectedResult,
			tracker.CompatibleWithTip(test.ComparedEvent),
			test.Description,
		)
	}
}
