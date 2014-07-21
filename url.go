package deje

import (
	"net/url"
	"strings"
)

func unbareUrl(deje_url string) string {
	valid_prefixes := []string{"http://", "https://", "ws://", "deje://", "//"}
	has_valid_prefix := false
	for _, prefix := range valid_prefixes {
		if strings.HasPrefix(deje_url, prefix) {
			has_valid_prefix = true
		}
	}
	if !has_valid_prefix {
		deje_url = "ws://" + deje_url
	}
	return deje_url
}

// Given a deje:// URL, return router URL and DEJE topic.
func GetRouterAndTopic(deje_url string) (router, topic string, err error) {
	var router_url, topic_url *url.URL
	router_url, err = url.Parse(unbareUrl(deje_url))
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
	topic_url.Scheme = "deje"
	if topic_url.Path == "" {
		topic_url.Path = "/"
	}

	return router_url.String(), topic_url.String(), nil
}
