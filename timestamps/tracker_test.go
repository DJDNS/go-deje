package timestamps

import (
	"bytes"
	"errors"
	"log"
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
type trackerDoIterationScenario struct {
	Description string
	Builder     trackerScenarioSetup
	Timestamps  []string
	Position    int
	Output      error
	StartTip    string
	EndTip      string
}
type trackerGoToLatestScenario struct {
	Description string
	Builder     trackerScenarioSetup
	Timestamps  []string
	LogOutput   string
	StartTip    string
	EndTip      string
}

func tsbuilderRaw() TimestampTracker {
	doc := document.NewDocument()
	service := NewPeerTimestampService(&doc)
	return NewTimestampTracker(&doc, service)
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

func TestTimestampTracker_GoToLatest(t *testing.T) {
	// For hash info
	dde := setupEvents(document.NewDocument())

	scenarios := []trackerGoToLatestScenario{
		trackerGoToLatestScenario{
			Description: "No timestamps, therefore no change",
			Builder:     tsbuilderRaw,
			Timestamps:  []string{},
			LogOutput:   "",
			StartTip:    "",
			EndTip:      "",
		},
		trackerGoToLatestScenario{
			Description: "StartIteration fails",
			Builder:     tsbuilderFails,
			Timestamps:  []string{},
			LogOutput:   "test_logger: tsbuilderFails() service breaks on purpose\n",
			StartTip:    "",
			EndTip:      "",
		},
		trackerGoToLatestScenario{
			Description: "StartIteration fails when already at a tip",
			Builder:     tsbuilderFails,
			Timestamps:  []string{},
			LogOutput:   "test_logger: tsbuilderFails() service breaks on purpose\n",
			StartTip:    "abc",
			EndTip:      "abc",
		},
		trackerGoToLatestScenario{
			Description: "Failure on one timestamp",
			Builder:     tsbuilderNormal,
			Timestamps:  []string{"xyz"},
			LogOutput:   "test_logger: Error on iteration 0 (current tip: ''):\ntest_logger: No such event\n",
			StartTip:    "",
			EndTip:      "",
		},
		trackerGoToLatestScenario{
			Description: "Success on one timestamp",
			Builder:     tsbuilderNormal,
			Timestamps:  []string{dde.Root.Hash()},
			LogOutput:   "",
			StartTip:    dde.Root.Hash(),
			EndTip:      dde.Root.Hash(),
		},
		trackerGoToLatestScenario{
			Description: "Failure does not impede or destroy progress",
			Builder:     tsbuilderNormal,
			Timestamps:  []string{"abc", dde.Root.Hash(), "xyz"},
			LogOutput:   "test_logger: Error on iteration 0 (current tip: ''):\ntest_logger: No such event\ntest_logger: Error on iteration 2 (current tip: '" + dde.Root.Hash() + "'):\ntest_logger: No such event\n",
			StartTip:    dde.Root.Hash(),
			EndTip:      dde.Root.Hash(),
		},
		trackerGoToLatestScenario{
			Description: "First fork wins",
			Builder:     tsbuilderNormal,
			Timestamps:  []string{dde.Root.Hash(), dde.Child.Hash(), dde.Fork.Hash()},
			LogOutput:   "test_logger: Error on iteration 2 (current tip: '" + dde.Child.Hash() + "'):\ntest_logger: Event is not compatible with and ahead of tip\n",
			StartTip:    dde.Child.Hash(),
			EndTip:      dde.Child.Hash(),
		},
	}
	for i, scenario := range scenarios {
		buf := new(bytes.Buffer)
		logger := log.New(buf, "test_logger: ", 0)

		tracker := scenario.Builder()
		tracker.Doc.Timestamps = scenario.Timestamps
		tracker.StartIteration()
		tracker.tip = scenario.StartTip

		assert.Equal(t,
			scenario.EndTip,
			tracker.GoToLatest(logger),
			"Scenario %d (%s)", i, scenario.Description,
		)
		assert.Equal(t, scenario.EndTip, tracker.tip)
		assert.Equal(t, scenario.LogOutput, buf.String())
	}
}
func TestTimestampTracker_GoToLatest_NilLogger(t *testing.T) {
	// Test StartIteration failure
	doc := document.NewDocument()
	service := failingTimestampService("tsbuilderFails() service breaks on purpose")
	tracker := NewTimestampTracker(&doc, service)

	assert.Equal(t, "", tracker.GoToLatest(nil))

	// Test DoIteration failure
	tracker.Service = NewPeerTimestampService(&doc)
	doc.Timestamps = []string{"This hash does not exist"}

	assert.Equal(t, "", tracker.GoToLatest(nil))
}

func TestTimestampTracker_DoIteration(t *testing.T) {
	// For hash info
	dde := setupEvents(document.NewDocument())

	scenarios := []trackerDoIterationScenario{
		trackerDoIterationScenario{
			Description: "No timestamps, therefore bad index",
			Builder:     tsbuilderRaw,
			Timestamps:  []string{},
			Position:    0,
			Output:      errors.New("Bad position"),
			StartTip:    "",
			EndTip:      "",
		},
		trackerDoIterationScenario{
			Description: "Timestamp references missing event",
			Builder:     tsbuilderNormal,
			Timestamps:  []string{dde.Unregistered.Hash()},
			Position:    0,
			Output:      errors.New("No such event"),
			StartTip:    "",
			EndTip:      "",
		},
		trackerDoIterationScenario{
			Description: "Incompatible branch",
			Builder:     tsbuilderNormal,
			Timestamps:  []string{dde.Root.Hash()},
			Position:    0,
			Output:      errors.New("Event is not compatible with and ahead of tip"),
			StartTip:    dde.Child.Hash(),
			EndTip:      dde.Child.Hash(),
		},
		trackerDoIterationScenario{
			Description: "Goto() failure",
			Builder:     tsbuilderNormal,
			Timestamps:  []string{dde.CannotGoto.Hash()},
			Position:    0,
			Output:      errors.New("No path provided"),
			StartTip:    "",
			EndTip:      "",
		},
		trackerDoIterationScenario{
			Description: "Success",
			Builder:     tsbuilderNormal,
			Timestamps:  []string{dde.Root.Hash(), dde.Child.Hash()},
			Position:    1,
			Output:      nil,
			StartTip:    dde.Child.Hash(),
			EndTip:      dde.Child.Hash(),
		},
	}
	for i, scenario := range scenarios {
		tracker := scenario.Builder()
		tracker.Doc.Timestamps = scenario.Timestamps
		tracker.StartIteration()
		tracker.tip = scenario.StartTip

		assert.Equal(t,
			scenario.Output,
			tracker.DoIteration(scenario.Position),
			"Scenario %d (%s)", i, scenario.Description,
		)
		assert.Equal(t, scenario.EndTip, tracker.tip)
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
