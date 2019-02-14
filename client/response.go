package client

import (
	"encoding/json"
	"net/http"
)

type Response func(target interface{}) error

func failResponse(err error) Response {
	return func(_ interface{}) error {
		return err
	}

}

func jsonDecodeResponse(resp *http.Response) Response {
	return func(target interface{}) error {
		// If the target is nil, just return nil to indicate no error was
		// encountered.
		if target == nil {
			return nil
		}

		defer resp.Body.Close()

		// Wrapper matches the standard response format from the Bifrost
		// server.
		wrapper := struct {
			Data interface{}
		}{target}

		return json.NewDecoder(resp.Body).Decode(&wrapper)
	}
}

// Returns the error encountered, if any, and doesn't try to parse the
// response body.
func (fn Response) Error() error {
	return fn(nil)
}

// Decode proxies directly to the Response function.
func (fn Response) Decode(target interface{}) error {
	return fn(target)
}
