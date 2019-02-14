package client

import (
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

	reduceOpts(opts...)(req)
	reduceOpts(
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
	return reduceOpts(
		Method("POST"),
		Path("channels"),
	)
}

func GetChannel(name string) FetchOpt {
	return reduceOpts(
		Method("GET"),
		Path("channels/"+name),
	)
}

//func (client *Client) ListChannels() ([]Channels, error) {
//	_, err := client.Fetch(Method("GET"), Path("channels"))
//	return nil, err
//}
//
//func (client *Client) DeleteChannel(name string) error {
//	_, err := client.Fetch(Method("DELETE"), Path("channels"))
//	return err
//}
