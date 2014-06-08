package util

type jsonObject map[string]interface{}

// No longer relevant for business logic, but fine for utils testing.
type ircLocation struct {
	Host    string `json:"host"`
	Port    uint32 `json:"port"`
	Channel string `json:"channel"`
}
