package middleware

import (
	"net/http"
)

func Wrap(handler http.Handler) (http.Handler, func(Middlewarer)) {
	app := &app{handler}

	// use will accept anything that quacks like middleware
	use := func(m Middlewarer) {
		m.Next(app.Handler)
		app.Handler = m
	}

	return app, use
}

type app struct {
	http.Handler
}

type Middlewarer interface {
	ServeHTTP(http.ResponseWriter, *http.Request)
	Next(http.Handler)
}

// Middleware can be embedded in an object to automatically implement the
// Middlewarer interface
type Middleware struct {
	http.Handler
}

// Sets the next handler for this middleware. This could be a
func (m *Middleware) Next(handler http.Handler) {
	m.Handler = handler
}
