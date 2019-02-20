package client

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"net/url"
	"strconv"
)

const apiBase = "/api/"

type Client struct {
	scheme string
	Host   string
	Port   uint

	cl http.Client
}

func New(host string, port uint) (*Client, error) {
	return &Client{
		scheme: "http",
		Host:   host,
		Port:   port,
	}, nil
}

func (client *Client) Do(opts ...FetchOpt) Response {
	req := &http.Request{
		URL: &url.URL{},
	}

	reduceOpts(
		reduceOpts(opts...),
		Proto("http"),
		Host(client.hostString()),
	)(req)

	response, err := client.cl.Do(req)
	if err != nil {
		return failResponse(err)
	}

	switch response.StatusCode {
	case 200:
		return jsonDecodeResponse(response)
	case 204:
		return emptyResponse
	default:
		return failHTTPStatus(response)
	}

	return jsonDecodeResponse(response)
}

func (client *Client) hostString() string {
	return client.Host + ":" + strconv.Itoa(int(client.Port))
}

func CreateChannel(name string, size uint) FetchOpt {
	buf, _ := json.Marshal(map[string]interface{}{
		"name": name,
		"size": size,
	})

	reader := ioutil.NopCloser(
		bytes.NewReader(buf),
	)

	return reduceOpts(
		Method("POST"),
		Path("channels"),
		Body(reader),
	)
}

func GetChannel(name string) FetchOpt {
	return reduceOpts(
		Method("GET"),
		Path("channels/"+name),
	)
}

func ListChannels(req *http.Request) {
	reduceOpts(
		Method("GET"),
		Path("channels"),
	)(req)
}

//func (client *Client) DeleteChannel(name string) error {
//	_, err := client.Fetch(Method("DELETE"), Path("channels"))
//	return err
//}
