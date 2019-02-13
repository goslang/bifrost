package handlers

import (
	"io/ioutil"
	"net/http"
	"time"

	"github.com/julienschmidt/httprouter"

	"github.com/goslang/bifrost/engine"
	"github.com/goslang/bifrost/server/lib/responder"
)

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
