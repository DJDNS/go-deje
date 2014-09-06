package timestamps

import (
	"errors"
	"testing"

	"github.com/DJDNS/go-deje/document"
	"github.com/stretchr/testify/assert"
)

type failingTimestampService string

func (fts failingTimestampService) GetTimestamps() ([]string, error) {
	return nil, errors.New(string(fts))
}

type demoDocEvents struct {
	Root         document.Event
	Child        document.Event
	Fork         document.Event
	Orphan       document.Event
	CannotGoto   document.Event
	Unregistered document.Event
}

func setupEvents(doc document.Document) demoDocEvents {
	var dde demoDocEvents

	dde.Root = doc.NewEvent("SET")
	dde.Root.Arguments["path"] = []interface{}{"key"}
	dde.Root.Arguments["value"] = "value"
	dde.Root.Register()

	dde.Child = doc.NewEvent("SET")
	dde.Child.Arguments["path"] = []interface{}{"other key"}
	dde.Child.Arguments["value"] = "other value"
	dde.Child.SetParent(dde.Root)
	dde.Child.Register()

	// Competes with dde.Child
	dde.Fork = doc.NewEvent("SET")
	dde.Fork.Arguments["path"] = []interface{}{"fork"}
	dde.Fork.Arguments["value"] = "fork"
	dde.Fork.SetParent(dde.Root)
	dde.Fork.Register()

	// No arguments
	dde.CannotGoto = doc.NewEvent("SET")
	dde.CannotGoto.Register()

	dde.Orphan = doc.NewEvent("SET")
	dde.Orphan.ParentHash = "foobarbaz"
	dde.Orphan.Register()

	dde.Unregistered = doc.NewEvent("foo")

	return dde
}

type trackerScenarioSetup func() TimestampTracker
type trackerFindLatestScenario struct {
	Description string
	Builder     trackerScenarioSetup
	Timestamps  []string
	Error       string
	TipHash     string
}

func tsbuilderNormal() TimestampTracker {
	doc := document.NewDocument()
	service := NewPeerTimestampService(&doc)

	setupEvents(doc)

	return NewTimestampTracker(&doc, service)
}
func tsbuilderFails() TimestampTracker {
	doc := document.NewDocument()
	service := failingTimestampService("tsbuilderFails() service breaks on purpose")
	return NewTimestampTracker(&doc, service)
}

func TestTimestampTracker_FindLatest(t *testing.T) {
	// For hash info
	dde := setupEvents(document.NewDocument())

	scenarios := []trackerFindLatestScenario{
		trackerFindLatestScenario{
			Description: "No timestamps, therefore no change",
			Builder:     tsbuilderNormal,
			Timestamps:  []string{},
			Error:       "",
			TipHash:     "",
		},
		trackerFindLatestScenario{
			Description: "GetTimestamps fails",
			Builder:     tsbuilderFails,
			Timestamps:  []string{},
			Error:       "tsbuilderFails() service breaks on purpose",
			TipHash:     "",
		},
		trackerFindLatestScenario{
			Description: "Failure on one timestamp",
			Builder:     tsbuilderNormal,
			Timestamps:  []string{"xyz"},
			Error:       "",
			TipHash:     "",
		},
		trackerFindLatestScenario{
			Description: "Success on one timestamp",
			Builder:     tsbuilderNormal,
			Timestamps:  []string{dde.Root.Hash()},
			Error:       "",
			TipHash:     dde.Root.Hash(),
		},
		trackerFindLatestScenario{
			Description: "Failure does not impede or destroy progress",
			Builder:     tsbuilderNormal,
			Timestamps:  []string{"abc", dde.Root.Hash(), "xyz"},
			Error:       "",
			TipHash:     dde.Root.Hash(),
		},
		trackerFindLatestScenario{
			Description: "First fork wins",
			Builder:     tsbuilderNormal,
			Timestamps:  []string{dde.Root.Hash(), dde.Child.Hash(), dde.Fork.Hash()},
			Error:       "",
			TipHash:     dde.Child.Hash(),
		},
	}
	for i, scenario := range scenarios {
		tracker := scenario.Builder()
		tracker.Doc.Timestamps = scenario.Timestamps

		event, err := tracker.FindLatest()
		var hash string
		if event != nil {
			hash = event.Hash()
		}

		assert.Equal(t, scenario.TipHash, hash, "Scenario %d (%s)", i, scenario.Description)
		if scenario.Error != "" {
			if assert.Error(t, err, scenario.Description) {
				assert.Equal(t, scenario.Error, err.Error())
			}
		} else {
			assert.NoError(t, err, scenario.Description)
		}
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
	dde := setupEvents(doc)

	tests := []compatibleTest{
		compatibleTest{
			Tip:            "",
			ComparedEvent:  nil,
			ExpectedResult: false,
			Description:    "A nil event pointer",
		},
		compatibleTest{
			Tip:            "",
			ComparedEvent:  &dde.Unregistered,
			ExpectedResult: false,
			Description:    "An unregistered event",
		},
		compatibleTest{
			Tip:            "",
			ComparedEvent:  &dde.Root,
			ExpectedResult: true,
			Description:    "Any registered event vs no-tip",
		},
		compatibleTest{
			Tip:            dde.Root.Hash(),
			ComparedEvent:  &dde.Child,
			ExpectedResult: true,
			Description:    "Child of root event",
		},
		compatibleTest{
			Tip:            dde.Child.Hash(),
			ComparedEvent:  &dde.Root,
			ExpectedResult: false,
			Description:    "Parent of tip event",
		},
		compatibleTest{
			Tip:            dde.Child.Hash(),
			ComparedEvent:  &dde.Fork,
			ExpectedResult: false,
			Description:    "Incompatible forks",
		},
		compatibleTest{
			Tip:            "foobar",
			ComparedEvent:  &dde.Root,
			ExpectedResult: false,
			Description:    "Some random broken tip value",
		},
		compatibleTest{
			Tip:            dde.Root.Hash(),
			ComparedEvent:  &dde.Orphan,
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
