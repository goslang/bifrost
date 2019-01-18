package handlers

import (
	"io/ioutil"
	"log"
	"net/http"

	"github.com/gorilla/websocket"
	"github.com/julienschmidt/httprouter"

	"github.com/goslang/bifrost/engine"
	"github.com/goslang/bifrost/lib/responder"
)

var upgrader = websocket.Upgrader{}

type channelController struct {
	engine *engine.Engine
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
	cc.engine.AddChannel(name)
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

func (cc *channelController) delete(
	w http.ResponseWriter,
	req *http.Request,
	p httprouter.Params,
) {
	name := p.ByName("name")
	cc.engine.RemoveChannel(name)

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
	defer c.Close()

	for {
		select {
		case message, ok := <-cc.engine.Listen(channelName):
			if !ok {
				break
			}

			err := c.WriteMessage(1, message)
			if err != nil {
				log.Println("write:", err)
				break
			}
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

	// NOTE: Ignoring error atm
	if err := cc.engine.Publish(name, buf); err != nil {
		log.Println("WARNING: Failed to publish message on channel", name)
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
	message, err := cc.engine.Pop(name)
	if err != nil {
		responder.New(w).Json(map[string]string{
			"error": err.Error(),
		})()
		return
	}

	responder.New(w).Json(map[string]string{
		"message": string(message),
	})()
}
