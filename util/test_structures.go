package util

type jsonObject map[string]interface{}

type ircLocation struct {
	Host    string `json:"host"`
	Port    uint32 `json:"port"`
	Channel string `json:"channel"`
}
