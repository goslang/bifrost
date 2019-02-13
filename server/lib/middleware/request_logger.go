package middleware

import (
	"log"
	"net/http"
)

type RequestLogger struct {
	Middleware
}

func (m *RequestLogger) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	contentType, ok := r.Header["Content-Type"]
	if !ok {
		contentType = []string{"-"}
	}

	log.Printf(`"%v %v" %v %db`,
		r.Method,
		r.URL,
		contentType,
		r.ContentLength,
	)

	m.Handler.ServeHTTP(w, r)
}
