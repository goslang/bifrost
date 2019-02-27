package handlers

import (
	"net/http"

	"github.com/gorilla/websocket"
	"github.com/julienschmidt/httprouter"

	"github.com/goslang/bifrost/engine"
	//"github.com/goslang/bifrost/server/lib/responder"
)

type subscriptionController struct {
	EventsCh chan engine.Event
	upgrader websocket.Upgrader
}

func NewSubscriptionController() *subscriptionController {
	return &subscriptionController{
		EventsCh: make(chan engine.Event),
		upgrader: websocket.Upgrader{
			CheckOrigin: func(_ *http.Request) bool { return true },
		},
	}
}

func (sc *subscriptionController) subscribe(
	w http.ResponseWriter,
	req *http.Request,
	p httprouter.Params,
) {
	//	channelName := p.ByName("name")
	//
	//	c, closed, err := sc.upgradeConnection(w, req, nil)
	//	if err != nil {
	//		responder.New(w).Status(400)()
	//	}
	//
	//	for {
	//		evt, messageCh := engine.Pop(req.Context(), channelName)
	//
	//		select {
	//		case sc.EventsCh <- evt:
	//		case <-closed:
	//			return
	//		}
	//
	//		select {
	//		case message, ok := <-messageCh:
	//			if !ok {
	//				return
	//			}
	//
	//			err := c.WriteMessage(websocket.TextMessage, message)
	//			if err != nil {
	//				return
	//			}
	//		case <-closed:
	//			return
	//		}
	//	}
}

func (sc *subscriptionController) upgradeConnection(
	w http.ResponseWriter,
	req *http.Request,
	responseHeader http.Header,
) (
	*websocket.Conn,
	chan struct{},
	error,
) {
	c, err := sc.upgrader.Upgrade(w, req, nil)
	if err != nil {
		return nil, nil, err
	}

	closed := make(chan struct{})
	go func() {
		defer func() {
			c.Close()
		}()

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

	return c, closed, nil
}
