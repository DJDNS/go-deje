package manager

import (
	"github.com/campadrenalin/go-deje/model"
	"reflect"
	"testing"
)

func TestTimestampManager_Register(t *testing.T) {
	m := NewTimestampManager()
	ts := model.Timestamp{
		QHash: "Bubba",
		Time:  25,
	}

	if m.Contains(ts) {
		t.Fatal("m should not contain ts yet")
	}
	m.Register(ts)
	if !m.Contains(ts) {
		t.Fatal("m should contain ts")
	}
}

func TestTimestampManager_Unregister(t *testing.T) {
	m := NewTimestampManager()
	ts := model.Timestamp{
		QHash: "Bubba",
		Time:  25,
	}
	m.Register(ts)

	if !m.Contains(ts) {
		t.Fatal("m should contain ts")
	}
	m.Unregister(ts)
	if m.Contains(ts) {
		t.Fatal("m should not contain ts anymore")
	}

	// Should be idempotent
	m.Unregister(ts)
	if m.Contains(ts) {
		t.Fatal("m should not contain ts anymore")
	}
}

func TestUint64Slice_Sort(t *testing.T) {
	uints := Uint64Slice{5, 0, 1000, 17}
	uints.Sort()

	expected := Uint64Slice{0, 5, 17, 1000}
	if !reflect.DeepEqual(uints, expected) {
		t.Fatalf("Expected %v, got %v", expected, uints)
	}
}

func setup_TSM_iter() (chan_ts, TimestampManager) {
	c := make(chan_ts)
	tm := NewTimestampManager()

	tm.Register(model.Timestamp{
		QHash: "hello",
		Time:  51,
	})
	tm.Register(model.Timestamp{
		QHash: "world",
		Time:  51,
	})
	tm.Register(model.Timestamp{
		QHash: "an early hash",
		Time:  51,
	})
	tm.Register(model.Timestamp{
		QHash: "an unrelated quorum",
		Time:  52,
	})

	return c, tm
}

func TestTimestampManager_emitBlock(t *testing.T) {
	c, tm := setup_TSM_iter()

	block := tm.GetGroup("51")
	go func() {
		tm.emitBlock(c, block)
		close(c)
	}()

	expected_qhashes := []string{"an early hash", "hello", "world"}
	i := 0
	for ts := range c {
		expected_qhash := expected_qhashes[i]
		i++
		if ts.QHash != expected_qhash {
			t.Fatalf("Expected %v, got %v", expected_qhash, ts.QHash)
		}
	}
	if i < len(expected_qhashes) {
		t.Fatalf("Expected %d results, got %d", len(expected_qhashes), i)
	}
}

func TestTimestampManager_emitTimestamps(t *testing.T) {
	c, tm := setup_TSM_iter()
	bh := Uint64Slice{3, 52, 51} // Crazy order, including empty blocks

	go func() {
		tm.emitTimestamps(c, bh)
		close(c)
	}()

	expected_qhashes := []string{"an unrelated quorum", "an early hash", "hello", "world"}
	i := 0
	for ts := range c {
		expected_qhash := expected_qhashes[i]
		i++
		if ts.QHash != expected_qhash {
			t.Fatalf("Expected %v, got %v", expected_qhash, ts.QHash)
		}
	}
	if i < len(expected_qhashes) {
		t.Fatalf("Expected %d results, got %d", len(expected_qhashes), i)
	}
}

func TestTimestampManager_sortedBlocks(t *testing.T) {
	_, tm := setup_TSM_iter()

	tm.GetGroup("4") // Also create empty group

	sorted_blocks, err := tm.sortedBlocks()
	if err != nil {
		t.Fatal(err)
	}
	expected := Uint64Slice{4, 51, 52}
	if !reflect.DeepEqual(sorted_blocks, expected) {
		t.Fatalf("Expected %v, got %v", expected, sorted_blocks)
	}

	tm.GetGroup("not a number") // Create invalid group
	sorted_blocks, err = tm.sortedBlocks()
	if err == nil {
		t.Fatal("tm.sortedBlocks() should have failed")
	}
}

func TestTimestampManager_Iter(t *testing.T) {
	_, tm := setup_TSM_iter()

	c, err := tm.Iter()
	if err != nil {
		t.Fatal(nil)
	}

	expected_qhashes := []string{"an early hash", "hello", "world", "an unrelated quorum"}
	i := 0
	for ts := range c {
		expected_qhash := expected_qhashes[i]
		i++
		if ts.QHash != expected_qhash {
			t.Fatalf("Expected %v, got %v", expected_qhash, ts.QHash)
		}
	}
	if i < len(expected_qhashes) {
		t.Fatalf("Expected %d results, got %d", len(expected_qhashes), i)
	}

	tm.GetGroup("not a number") // Create invalid group
	_, err = tm.Iter()
	if err == nil {
		t.Fatal("tm.Iter() should have failed")
	}
}
