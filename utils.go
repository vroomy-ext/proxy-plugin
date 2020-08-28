package main

import (
	"fmt"
	"net/http"
	"net/http/httputil"
	"net/url"

	"github.com/Hatch1fy/httpserve"
	"github.com/hatchify/atoms"
)

func getURLs(args []string) (urls []string, err error) {
	urls = args
	return
}

func setHost(req *http.Request) {
	fmt.Println("Setting!", host)
	req.Host = host
	req.URL.Host = host
	req.URL.Scheme = "http"
	fmt.Println("Set!!", req)
}

func newHandler(u *url.URL) httpserve.Handler {
	handler := httputil.NewSingleHostReverseProxy(u)

	// Wrap gofast.Handler with httpserve Handler
	return func(ctx *httpserve.Context) (res httpserve.Response) {
		// Call handler.ServeHTTP and pass it the writer and request
		fmt.Println("Request!", ctx.Request)
		handler.ServeHTTP(ctx.Writer, ctx.Request)
		return
	}
}

func newHandlersFromArgs(args []string) (hs []httpserve.Handler, err error) {
	var urls []string
	// Get the URL from the provided arguments
	if urls, err = getURLs(args); err != nil {
		return
	}

	for _, urlStr := range urls {
		var u *url.URL
		if u, err = url.Parse(urlStr); err != nil {
			return
		}

		hs = append(hs, newHandler(u))
	}

	return
}

func newRoundRobinHandler(hs []httpserve.Handler) (h httpserve.Handler) {
	var i atoms.Uint64
	// Set handlers length
	hsl := uint64(len(hs))
	// Set index to last position so the first load is the first Handler
	i.Store(hsl - 1)

	return func(ctx *httpserve.Context) (res httpserve.Response) {
		// Set value as an atomic increment of index
		val := i.Add(1)
		// Get the index from the value
		index := val % hsl
		// Grab Handler from Handlers slice
		handler := hs[index]
		// Return retrieved Handler
		return handler(ctx)
	}
}
