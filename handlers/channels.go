package handlers

import (
	"io/ioutil"
	"net/http"
	"time"

	"github.com/gorilla/websocket"
	"github.com/julienschmidt/httprouter"

	"github.com/goslang/bifrost/engine"
	"github.com/goslang/bifrost/lib/responder"
)

var upgrader = websocket.Upgrader{
	CheckOrigin: func(_ *http.Request) bool {
		// Allow everything right now. Our CORS Middleware should protect us
		// from any attacks.
		return true
	},
}

type channelController struct {
	engine   *engine.Engine
	eventsCh chan engine.Event
}

func NewChannelController(eng *engine.Engine) *channelController {
	return &channelController{
		engine: eng,
	}
}

func (cc *channelController) create(
	w http.ResponseWriter,
	req *http.Request,
	p httprouter.Params,
) {
	name := p.ByName("name")
	cc.eventsCh <- engine.AddChannel(name)

	responder.New(w).Json(map[string]string{
		"status": "ok",
	})()
}

func (cc *channelController) get(
	w http.ResponseWriter,
	req *http.Request,
	_ httprouter.Params,
) {
	responder.New(w).Json(map[string]string{
		"status": "ok",
	})()
}

func (cc *channelController) list(
	w http.ResponseWriter,
	req *http.Request,
	p httprouter.Params,
) {
	responder.New(w).Json(map[string]string{
		"status": "ok",
	})()
}

func (cc *channelController) destroy(
	w http.ResponseWriter,
	req *http.Request,
	p httprouter.Params,
) {
	name := p.ByName("name")
	cc.eventsCh <- engine.RemoveChannel(name)

	responder.New(w).Json(map[string]string{
		"status": "ok",
	})()
}

func (cc *channelController) subscribe(
	w http.ResponseWriter,
	req *http.Request,
	p httprouter.Params,
) {
	channelName := p.ByName("name")

	c, err := upgrader.Upgrade(w, req, nil)
	if err != nil {
		return
	}
	defer func() {
		c.Close()
	}()

	closed := make(chan bool)
	go func() {
		for {
			// Watch for a closed connection. Apparently writing to a closed
			// connection doesn't fail as expected, so continually try to read
			// and if that fails, signal the main loop to exit.
			_, _, err := c.ReadMessage()
			if err != nil {
				close(closed)
				return
			}
		}
	}()

	listener, messageCh := engine.ChangeListener(channelName)
	// defer close(messageCh) // Not safe, other goroutine maybe writing.

	id := cc.engine.Register(listener)
	defer cc.engine.Deregister(id)

	for {
		select {
		case message, ok := <-messageCh:
			if !ok {
				return
			}

			err := c.WriteMessage(websocket.TextMessage, message)
			if err != nil {
				return
			}
		case <-closed:
			return
		}
	}
}

func (cc *channelController) publish(
	w http.ResponseWriter,
	req *http.Request,
	p httprouter.Params,
) {
	name := p.ByName("name")

	buf, err := ioutil.ReadAll(req.Body)
	if err != nil {
		responder.New(w).
			Status(http.StatusBadRequest).
			Json(map[string]string{
				"error": "Failed to read request body",
			})()
	}

	cc.eventsCh <- &engine.PushEvent{
		QueueName: name,
		Data:      buf,
	}

	responder.New(w).Json(map[string]string{
		"status": "ok",
	})()
}

func (cc *channelController) pop(
	w http.ResponseWriter,
	req *http.Request,
	p httprouter.Params,
) {
	name := p.ByName("name")
	ch := make(chan []byte)
	defer close(ch)

	cc.eventsCh <- engine.Pop(name, func(message []byte) {
		ch <- message
	})

	select {
	case message := <-ch:
		responder.New(w).Json(map[string]string{
			"message": string(message),
		})()
	case <-time.After(10 * time.Second):
		responder.New(w).
			Status(http.StatusServiceUnavailable).
			Json(map[string]string{})()
	}
}
