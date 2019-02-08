package responder

import (
	"encoding/json"
	"net/http"
	"time"
)

type Responder func() (http.ResponseWriter, error)

func New(w http.ResponseWriter) Responder {
	return func() (http.ResponseWriter, error) { return w, nil }
}

func (r Responder) Headers(headers map[string][]string) Responder {
	return func() (w http.ResponseWriter, err error) {
		defer errorHandler(&err)
		w, err = r()
		checkErr(err)
		finalHeaders := w.Header()
		for key, value := range headers {
			if header, exists := finalHeaders[key]; exists {
				finalHeaders[key] = append(header, value...)
			} else {
				finalHeaders[key] = value
			}
		}
		return
	}
}

func (r Responder) Status(status int) Responder {
	return func() (w http.ResponseWriter, err error) {
		defer errorHandler(&err)
		w, err = r()
		checkErr(err)
		w.WriteHeader(status)
		return
	}
}

func (r Responder) Json(data interface{}) Responder {
	return func() (w http.ResponseWriter, err error) {
		defer errorHandler(&err)
		w, err = r()
		checkErr(err)
		jsonB, err := json.Marshal(data)
		checkErr(err)
		w.Write([]byte(jsonB))
		return
	}
}

func (r Responder) Cookie(key, value, path string, secure bool) Responder {
	return func() (w http.ResponseWriter, err error) {
		defer errorHandler(&err)
		w, err = r()
		checkErr(err)
		http.SetCookie(w, &http.Cookie{
			Name:     key,
			Value:    value,
			Path:     path,
			Secure:   secure,
			HttpOnly: true,
			Expires:  time.Now().Add(2 * time.Hour * 24 * 30 * 12),
		})
		return
	}
}

func (r Responder) Must() http.ResponseWriter {
	w, err := r()
	checkErr(err)
	return w
}

func errorHandler(errP *error) {
	if err := recover(); err != nil {
		*errP = err.(error)
	}
}

func checkErr(err error) {
	if err != nil {
		panic(err)
	}
}
