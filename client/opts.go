package client

import (
	"io"
	"net/http"
)

type FetchOpt func(*http.Request)

func reduceOpts(opts ...FetchOpt) FetchOpt {
	return func(req *http.Request) {
		for _, opt := range opts {
			opt(req)
		}
	}
}

func Method(method string) FetchOpt {
	return func(req *http.Request) {
		req.Method = method
	}
}

func Path(path string) FetchOpt {
	return func(req *http.Request) {
		req.URL.Path = apiBase + path
	}
}

func Body(reader io.ReadCloser) FetchOpt {
	return func(req *http.Request) {
		req.Body = reader
	}
}

func Host(host string) FetchOpt {
	return func(req *http.Request) {
		req.URL.Host = host
	}
}

func Proto(proto string) FetchOpt {
	return func(req *http.Request) {
		req.URL.Scheme = proto
	}
}
