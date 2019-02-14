package client

import (
	"fmt"

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

//func (client *Client) Fetch(opts ...FetchOpt) Response {
//	req := applyOpts(&http.Request{})
//
//	_, err := client.cl.Do(req)
//	return failResponse(err)
//}

func (client *Client) Do(opts ...FetchOpt) Response {
	req := &http.Request{
		URL: &url.URL{},
	}

	reduceOpts(opts...)(req)
	reduceOpts(
		Proto("http"),
		Host(client.hostString()),
	)(req)

	fmt.Println(req)
	response, err := client.cl.Do(req)
	if err != nil {
		return failResponse(err)
	}

	return jsonDecodeResponse(response)
}

//func (client *Client) buildURL(path string) url.URL {
//	return url.URL{
//		Host: client.Host + ":" + strconv.Itoa(int(client.Port)),
//		Path: apiBase + path,
//	}
//}

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
