package client

import (
	"bytes"
	"io/ioutil"
	"net/http"
)

type Channel struct {
	Name string
	Size uint
	Max  uint
}

func NewChannel(name string) *Channel {
	return &Channel{
		Name: name,
	}
}

func (c *Channel) Publish(message string) FetchOpt {
	body := ioutil.NopCloser(
		bytes.NewReader([]byte(message)),
	)

	return reduceOpts(
		Method("POST"),
		Path("channels/"+c.Name+"/publish"),
		Body(body),
	)
}

func (c *Channel) Pop(req *http.Request) {
	reduceOpts(
		Method("PATCH"),
		Path("channels/"+c.Name+"/pop"),
	)(req)
}
