package timestamps

import "testing"

func TestDTS_GetTimestamps(t *testing.T) {
	dts := DummyTimestampService{}
	stamps, err := dts.GetTimestamps("Interstella")
	if len(stamps) != 0 {
		t.Fatal(
			"Expected empty timestamp array, has %d elements",
			len(stamps),
		)
	}
	if err != nil {
		t.Fatal(err)
	}
}
