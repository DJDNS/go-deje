package deje

import (
	"encoding/json"
	"errors"
)

type JSONObject map[string]interface{}
type DocumentState JSONObject

type UnsetFieldError struct {
	FieldName string
}

func (e *UnsetFieldError) Error() string {
	return "Field " + e.FieldName + " was not set."
}

func FillStruct(m JSONObject, s interface{}) error {
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

func (s DocumentState) GetChannel() (*IRCLocation, error) {
	data, ok := s["channel"]
	if !ok {
		return nil, errors.New("Document does not have channel data")
	}
	m, ok := data.(JSONObject)
	if !ok {
		return nil, errors.New("Channel data was wrong type")
	}

	channel := new(IRCLocation)
	err := FillStruct(m, channel)
	return channel, err
}
