package model

import (
	"net"
	"net/url"
	"strconv"
)

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

// Convenience function to get data as URL variables
func (loc IRCLocation) GetVariables() string {
	str_port := strconv.FormatUint(uint64(loc.Port), 10)
	return "host=" + loc.Host +
		"&port=" + str_port +
		"&channel=" + loc.Channel
}

func (loc *IRCLocation) ParseFrom(urlstr string) error {
	urlobj, err := url.Parse(urlstr)
	if err != nil {
		return err
	}

	host, port, err := net.SplitHostPort(urlobj.Host)
	if err != nil {
		host, port = urlobj.Host, "6667"
	}
	portint, err := strconv.ParseUint(port, 10, 32)
	if err != nil {
		return err
	}

	loc.Host = host
	loc.Port = uint32(portint)
	loc.Channel = urlobj.Fragment

	return nil
}
