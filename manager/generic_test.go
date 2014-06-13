package manager

import (
	"github.com/campadrenalin/go-deje/model"
	"testing"
)

func setup_quorums() (model.Quorum, model.Quorum) {
	A := model.Quorum{
		EventHash:  "xyz",
		Signatures: map[string]string{"foo": "bar"},
	}
	B := model.Quorum{
		EventHash:  "xyz",
		Signatures: map[string]string{"fizz": "buzz"},
	}
	return A, B
}

func setup_om_with_ab() (genericManager, model.Quorum, model.Quorum) {
	m := newGenericManager()

	A, B := setup_quorums()

	m.register(A)
	m.register(B)

	return m, A, B
}

func TestGenericManagerGetItems(t *testing.T) {
	m, A, B := setup_om_with_ab()

	items := m.GetItems()
	if len(items) != 2 {
		t.Fatalf("Expected 2 items, got %d", len(items))
	}

	for _, q := range []model.Quorum{A, B} {
		key := q.GetKey()
		item, ok := items[key]
		if !ok {
			t.Fatalf("Missing item %s", key)
		}
		if !q.Eq(item) {
			t.Fatalf("TS %#v does not equal %#v", q, item)
		}
	}
}

func TestGenericManagerGetByKey(t *testing.T) {
	m, A, B := setup_om_with_ab()

	for _, q := range []model.Quorum{A, B} {
		key := q.GetKey()
		item, ok := m.GetByKey(key)
		if !ok {
			t.Fatalf("Missing item %s", key)
		}
		if !q.Eq(item) {
			t.Fatalf("TS %#v does not equal %#v", q, item)
		}
	}
}

func TestGenericManagerRegister(t *testing.T) {
	m, A, B := setup_om_with_ab()

	group := m.GetGroup("xyz")
	if !(group.Contains(A) && group.Contains(B)) {
		t.Fatal("group was missing quorums")
	}
}

func TestGenericManagerUnregister(t *testing.T) {
	m, A, B := setup_om_with_ab()

	m.unregister(A)

	group := m.GetGroup("xyz")
	if group.Contains(A) || !group.Contains(B) || m.Length() != 1 {
		t.Fatal("Failed to unregister A correctly")
	}
}

func TestManagableSetContains(t *testing.T) {
	ms := make(model.ManageableSet)
	A, B := setup_quorums()

	if ms.Contains(A) {
		t.Fatal("ms shouldn't contain A")
	}

	ms[A.GetKey()] = B
	if ms.Contains(A) {
		t.Fatal("ms shouldn't contain A, but has key for A")
	}

	ms[A.GetKey()] = A
	if !ms.Contains(A) {
		t.Fatal("ms should contain A")
	}
}

func TestGenericManagerContains(t *testing.T) {
	m := newGenericManager()
	A, B := setup_quorums()

	m.register(A)

	if !m.Contains(A) {
		t.Fatal("m should contain A")
	}
	if m.Contains(B) {
		t.Fatal("m should not contain B")
	}
}

func TestGenericManagerGetGroup(t *testing.T) {
	m := newGenericManager()
	group := m.GetGroup("xyz")

	if len(group) != 0 {
		t.Fatal("group should have been empty")
	}

	quorum, _ := setup_quorums()
	m.register(quorum)

	group = m.GetGroup("abc")
	if len(group) != 0 {
		t.Fatal("group should still have been empty")
	}
	group = m.GetGroup("xyz")
	if len(group) != 1 {
		t.Fatal("group should not have been empty")
	}

	if !quorum.Eq(group[quorum.GetKey()]) {
		t.Fatalf("group %#v should have contained %#v", group, quorum)
	}

	if len(m.by_group) != 2 {
		t.Fatalf(
			"Not caching groups - expected 2 groups, got %d: %#v",
			len(m.by_group),
			m.by_group,
		)
	}
}
