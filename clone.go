package deje

import "encoding/json"

type JSONObject map[string]interface{}

type UnsetFieldError struct {
	FieldName string
}

func (e *UnsetFieldError) Error() string {
	return "Field " + e.FieldName + " was not set."
}

func CloneMarshal(m interface{}, s interface{}) error {
	// Pass in &YourStructType{}
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
