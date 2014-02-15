package manager

import (
	"github.com/campadrenalin/go-deje/model"
	"testing"
)

func TestQuorumManager_Register(t *testing.T) {
	m := NewQuorumManager()
	q := model.NewQuorum("evhash")

	if m.Contains(q) {
		t.Fatal("m should not contain q yet")
	}
	m.Register(q)
	if !m.Contains(q) {
		t.Fatal("m should contain q")
	}
}

func TestQuorumManager_Unregister(t *testing.T) {
	m := NewQuorumManager()
	q := model.NewQuorum("evhash")
	m.Register(q)

	if !m.Contains(q) {
		t.Fatal("m should contain q")
	}
	m.Unregister(q)
	if m.Contains(q) {
		t.Fatal("m should not contain q anymore")
	}

	// Should be idempotent
	m.Unregister(q)
	if m.Contains(q) {
		t.Fatal("m should not contain q anymore")
	}
}

func TestQuorumManager_DeserializeFrom(t *testing.T) {
	m := NewQuorumManager()
	serial := make(map[string]model.Quorum)

	// Use fake keys - should be ignored
	serial["energy drink"] = model.NewQuorum("on the ground")
	serial["hot dog"] = model.NewQuorum("hand outs")
	serial["cell phone"] = model.NewQuorum("not my dad")

	m.DeserializeFrom(serial)

	items := m.GetItems()
	for key, value := range serial {
		_, ok := items[key]
		if ok {
			t.Fatalf("Key %s should not have been present in items", key)
		}

		if !m.Contains(value) {
			t.Fatalf("m should contain %#v", value)
		}
	}
}

func TestQuorumManager_SerializeTo(t *testing.T) {
	m := NewQuorumManager()
	serial := make(map[string]model.Quorum)

	quorums := []model.Quorum{
		model.NewQuorum("A"),
		model.NewQuorum("B"),
		model.NewQuorum("C"),
	}
	for _, value := range quorums {
		m.Register(value)
	}

	m.SerializeTo(serial)

	if len(serial) != len(quorums) {
		t.Fatalf("Expected %d quorums, got %d", len(quorums), len(serial))
	}
	for _, value := range quorums {
		key := value.GetKey()
		serial_quorum, ok := serial[key]
		if !ok {
			t.Fatalf("Item %s not present", key)
		}
		if !value.Eq(serial_quorum) {
			t.Fatalf("%#v != %#v", value, serial_quorum)
		}
	}
}
