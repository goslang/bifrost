package middleware

import (
	"net/http"
	"strings"
)

type Cors struct {
	Origins set
	Middleware

	Headers     []string
	Methods     []string
	Credentials bool
}

type set map[string]bool

var headers []string = []string{"Content-Type", "Authorization"}
var methods []string = []string{"GET", "POST", "DELETE", "PATCH", "PUT"}
var bools map[bool]string = map[bool]string{true: "true", false: "false"}

func NewCors(origins ...string) *Cors {
	cors := &Cors{}
	cors.Origins = make(set)

	for _, origin := range origins {
		cors.Origins[origin] = true
	}

	cors.Headers = append(cors.Headers, headers...)
	cors.Methods = append(cors.Headers, methods...)
	cors.Credentials = true

	return cors
}

func (m *Cors) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	headers := w.Header()
	origin := r.Header.Get("Origin")

	if m.originIsAllowed(origin) {
		headers["Access-Control-Allow-Origin"] = []string{origin}
		headers["Access-Control-Allow-Headers"] = []string{strings.Join(m.Headers, ", ")}
		headers["Access-Control-Allow-Methods"] = []string{strings.Join(m.Methods, ", ")}
		headers["Access-Control-Allow-Credentials"] = []string{bools[m.Credentials]}
	}

	m.Handler.ServeHTTP(w, r)
}

func (m *Cors) originIsAllowed(origin string) bool {
	_, included := m.Origins[origin]
	return included
}
