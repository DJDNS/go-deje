package deje

import (
	"github.com/campadrenalin/go-deje/model"
	"testing"
)

func TestGetDocument(t *testing.T) {
	con := NewDEJEController()
	loc := model.IRCLocation{"hello", 1234, "world"}
	doc := con.GetDocument(loc)

	if doc.Channel != loc {
		t.Fatalf("Expected %v, got %v", loc, doc.Channel)
	}
}
