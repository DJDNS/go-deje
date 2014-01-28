package deje

import "errors"

type DocumentState JSONObject

func (ds DocumentState) GetProperty(name string, s interface{}) error {
	data, ok := ds[name]
	if !ok {
		return errors.New("Document does not have "+name+" property")
	}

    return CloneMarshal(data, s)
}

func (ds DocumentState) GetChannel() (*IRCLocation, error) {
	channel := new(IRCLocation)
    err := ds.GetProperty("channel", channel)
	return channel, err
}

func (s DocumentState) GetURLs() (*DownloadURLs, error) {
	urls := new(DownloadURLs)
    err := s.GetProperty("urls", urls)
	return urls, err
}
