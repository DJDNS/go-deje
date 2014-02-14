package deje

import (
	"github.com/campadrenalin/go-deje/model"
	"testing"
)

func TestGetDocument(t *testing.T) {
	con := DEJEController{}
	con.GetDocument(model.IRCLocation{})
}
