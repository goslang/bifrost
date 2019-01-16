package handlers

import (
	"net/http"

	"github.com/julienschmidt/httprouter"
)

type channelController struct {
	list      http.Handler
	get       http.Handler
	create    http.Handler
	subscribe http.Handler
}

func NewRouter() *httprouter.Router {
	r := httprouter.New()

	r.GET("/api/channels", listChannels)
	r.GET("/api/channels/:name", getChannel)
	r.POST("/api/channels/:name", createChannel)

	r.GET("/api/channels/:name/subscribe", subscribe)

	return r
}
