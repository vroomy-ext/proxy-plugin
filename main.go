package main

import (
	"net/url"

	"github.com/Hatch1fy/httpserve"
	"github.com/hatchify/errors"
)

const (
	// ErrEmptyURL is returned when a URL is empty
	ErrEmptyURL = errors.Error("URL is empty")
)

var (
	host string
)

// Init is called when Vroomy initializes the plugin
func Init(env map[string]string) (err error) {
	host = env["proxy-host"]
	return
}

// Proxy will perform a direct proxy for a single URL
func Proxy(args ...string) (h httpserve.Handler, err error) {
	if len(args) == 0 {
		err = ErrEmptyURL
		return
	}

	var u *url.URL
	if u, err = url.Parse(args[0]); err != nil {
		return
	}

	h = newHandler(u)
	return
}

// RoundRobin will perform a round robin proxy for each of the provided URLs
func RoundRobin(urls ...string) (h httpserve.Handler, err error) {
	var hs []httpserve.Handler
	if hs, err = newHandlersFromArgs(urls); err != nil {
		return
	}

	h = newRoundRobinHandler(hs)
	return
}
