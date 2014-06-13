package services

import "testing"

func TestDTS_GetAfter(t *testing.T) {
	dts := DummyTimestampService{}
	dts.GetAfter("Interstella", 5555)
}

func TestDTS_MakeTimestamp(t *testing.T) {
	dts := DummyTimestampService{}
	dts.MakeTimestamp("hello", "world")
}
