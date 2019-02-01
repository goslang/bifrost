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
	EventsCh chan engine.Event
}

func NewChannelController() *channelController {
	return &channelController{
		EventsCh: make(chan engine.Event),
	}
}

func (cc *channelController) create(
	w http.ResponseWriter,
	req *http.Request,
	p httprouter.Params,
) {
	name := p.ByName("name")
	cc.EventsCh <- engine.AddChannel(name, 5)

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
	cc.EventsCh <- engine.RemoveChannel(name)

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

	for {
		evt, messageCh := engine.Pop(req.Context(), channelName)

		select {
		case cc.EventsCh <- evt:
		case <-closed:
			return
		}

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

	evt, confirmCh := engine.PushMessage(name, buf)
	cc.EventsCh <- evt

	_, ok := <-confirmCh
	if !ok {
		responder.New(w).Json(map[string]string{
			"status": "error",
		})()
		return
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

	evt, ch := engine.PopNow(name)
	cc.EventsCh <- evt

	select {
	case message, ok := <-ch:
		if !ok {
			responder.New(w).Status(404).Json(map[string]string{
				"status": "error",
			})()
			return
		}

		responder.New(w).Json(map[string]string{
			"message": string(message),
		})()
	case <-time.After(5 * time.Second):
		responder.New(w).
			Status(http.StatusServiceUnavailable).
			Json(map[string]string{
				"error": "timeout",
			})()
	}
}
