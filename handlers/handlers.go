package handlers

import (
	"github.com/julienschmidt/httprouter"

	"github.com/goslang/bifrost/engine"
)

func NewRouter(eventsCh chan<- engine.Event) *httprouter.Router {
	r := httprouter.New()

	channels := NewChannelController(eventsCh)

	r.GET("/api/channels", channels.list)
	r.GET("/api/channels/:name", channels.get)
	r.POST("/api/channels/:name", channels.create)
	r.DELETE("/api/channels/:name", channels.destroy)

	r.GET("/api/channels/:name/subscribe", channels.subscribe)
	r.POST("/api/channels/:name/publish", channels.publish)
	r.PATCH("/api/channels/:name/pop", channels.pop)

	return r
}
