package model

import "strconv"

// Describes an IRC Server+Channel combo. Every DEJE doc has an
// IRC location for broadcast communication.
type IRCLocation struct {
	Host    string `json:"host"`
	Port    uint32 `json:"port"`
	Channel string `json:"channel"`
}

func (loc IRCLocation) GetURL() string {
	str_port := strconv.FormatUint(uint64(loc.Port), 10)
	return "irc://" + loc.Host + ":" + str_port + "/#" + loc.Channel
}
