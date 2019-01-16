package handlers

import (
	"log"
	"net/http"

	"github.com/gorilla/websocket"
	"github.com/julienschmidt/httprouter"

	"github.com/goslang/bifrost/lib/responder"
)

var upgrader = websocket.Upgrader{}

// magicMap will be removed, but for now it maps channel names to go-channels.
var magicMap = make(map[string]chan []byte)

func createChannel(w http.ResponseWriter, req *http.Request, p httprouter.Params) {
	channelName := p.ByName("name")

	magicMap[channelName] = make(chan []byte)
}

func getChannel(w http.ResponseWriter, req *http.Request, _ httprouter.Params) {
	responder.New(w).Json(map[string]string{"status": "ok"}).Must()
}

func listChannels(w http.ResponseWriter, req *http.Request, p httprouter.Params) {
}

func subscribe(w http.ResponseWriter, req *http.Request, p httprouter.Params) {
	channelName := p.ByName("name")

	// TODO: Move to error handler somewhere else
	ch, ok := magicMap[channelName]
	if !ok {
		responder.New(w).
			Status(404).
			Json(map[string]string{
				"error": "Not found",
			}).Must()
	}

	c, err := upgrader.Upgrade(w, req, nil)
	if err != nil {
		return
	}
	defer c.Close()

	for {
		message := <-ch

		err := c.WriteMessage(1, message)
		if err != nil {
			log.Println("write:", err)
			break
		}
	}
}
