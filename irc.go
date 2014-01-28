package deje

type IRCLocation struct {
	Host    string `json:"host"`
	Port    uint32 `json:"port"`
	Channel string `json:"channel"`
}
