package deje

import (
	"errors"
	"net/url"
	"strings"
)

// Given a deje:// URL, return router URL and DEJE topic.
func GetRouterAndTopic(deje_url string) (router, topic string, err error) {
	var router_url, topic_url *url.URL
	if !strings.HasPrefix(deje_url, "deje://") {
		return "", "", errors.New("URL does not start with 'deje://': '" + deje_url + "'")
	}

	router_url, err = url.Parse(deje_url)
	if err != nil {
		return "", "", err
	}

	// Separate copy, each gets different manipulations
	topic_url = new(url.URL)
	*topic_url = *router_url

	// Router URL manipulations
	router_url.Scheme = "ws"
	router_url.Path = "/ws"

	// Topic URL manipulations
	if topic_url.Path == "" {
		topic_url.Path = "/"
	}

	return router_url.String(), topic_url.String(), nil
}
