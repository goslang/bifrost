package handlers

import (
	"context"

	"github.com/julienschmidt/httprouter"

	"github.com/goslang/bifrost/engine"
)

func NewRouter(eng *engine.Engine) *httprouter.Router {
	r := httprouter.New()

	channels := NewChannelController()
	subscriptions := NewSubscriptionController()

	r.GET("/api/channels", channels.list)
	r.GET("/api/channels/:name", channels.get)
	r.POST("/api/channels/:name", channels.create)
	r.DELETE("/api/channels/:name", channels.destroy)

	r.POST("/api/channels/:name/publish", channels.publish)
	r.PATCH("/api/channels/:name/pop", channels.pop)

	r.GET("/api/subscribe/:name", subscriptions.subscribe)

	ctx := context.Background()
	eng.Process(ctx, channels.EventsCh)
	eng.Process(ctx, subscriptions.EventsCh)

	return r
}
