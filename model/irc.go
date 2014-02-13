package model

// Describes an IRC Server+Channel combo. Every DEJE doc has an
// IRC location for broadcast communication.
type IRCLocation struct {
	Host    string `json:"host"`
	Port    uint32 `json:"port"`
	Channel string `json:"channel"`
}
