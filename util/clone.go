package util

import "encoding/json"

// Uses json.Marshal and json.Unmarshal to copy from
// one object to another. Excellent way to turn a JSON Object
// into a struct, or vice-versa.
//
// This makes its changes in-place, so always pass in the
// pointer to the struct object you're trying to fill.
func CloneMarshal(m interface{}, s interface{}) error {
	jstr, err := json.Marshal(m)
	if err != nil {
		return err
	}

	err = json.Unmarshal(jstr, s)
	if err != nil {
		return err
	}

	return nil
}
