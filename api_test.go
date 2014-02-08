package deje

import (
	"github.com/campadrenalin/go-deje/serial"
	"testing"
)

func TestGetDocument(t *testing.T) {
	con := DEJEController{}
	con.GetDocument(serial.IRCLocation{})
}
