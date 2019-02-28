package handlers

import (
	//"encoding/json"
	"io/ioutil"
	"net/http"
	"time"

	"github.com/julienschmidt/httprouter"
	"github.com/satori/go.uuid"

	"github.com/goslang/bifrost/engine"
	"github.com/goslang/bifrost/server/lib/responder"
)

type channelController struct {
	Stats       engine.StatsAPI
	ListenerAPI engine.ListenerAPI
	EventsCh    chan engine.Event
}

func NewChannelController(
	stats engine.StatsAPI,
	listenerAPI engine.ListenerAPI,
) *channelController {
	return &channelController{
		Stats:       stats,
		ListenerAPI: listenerAPI,
		EventsCh:    make(chan engine.Event),
	}
}

func (cc *channelController) create(
	w http.ResponseWriter,
	req *http.Request,
	_ httprouter.Params,
) {
	//	var parsed struct {
	//		Name string
	//		Size uint
	//	}
	//
	//	if err := json.NewDecoder(req.Body).Decode(&parsed); err != nil {
	//		responder.New(w).
	//			Status(http.StatusUnprocessableEntity).
	//			Json(map[string]string{
	//				"status": "error",
	//				"Error":  "Unprocessable Entity",
	//			})()
	//		return
	//	}
	//
	//	cc.EventsCh <- engine.AddChannel(parsed.Name, parsed.Size)
	//
	//	responder.New(w).Json(map[string]string{
	//		"status": "ok",
	//	})()
}

func (cc *channelController) get(
	w http.ResponseWriter,
	req *http.Request,
	p httprouter.Params,
) {
	name := p.ByName("name")

	detail, ok := cc.Stats.GetQueueDetails(name)
	if !ok {
		responder.New(w).
			Status(http.StatusNotFound).
			Json(map[string]interface{}{
				"status": "not found",
			})()
		return
	}

	responder.New(w).Json(map[string]interface{}{
		"status": "ok",
		"data":   detail,
	})()
}

func (cc *channelController) list(
	w http.ResponseWriter,
	req *http.Request,
	p httprouter.Params,
) {
	queues := cc.Stats.ListQueues()

	responder.New(w).Json(map[string]interface{}{
		"status": "ok",
		"data":   queues,
	})()
}

func (cc *channelController) destroy(
	w http.ResponseWriter,
	req *http.Request,
	p httprouter.Params,
) {
	//	name := p.ByName("name")
	//	cc.EventsCh <- engine.RemoveChannel(name)
	//
	//	responder.New(w).Json(map[string]string{
	//		"status": "ok",
	//	})()
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

	cc.EventsCh <- engine.Push{
		Channel: name,
		Message: buf,
	}

	//	evt, confirmCh := engine.PushMessage(name, buf)
	//	cc.EventsCh <- evt
	//
	//	_, ok := <-confirmCh
	//	if !ok {
	//		responder.New(w).
	//			Status(http.StatusConflict).
	//			Json(map[string]string{
	//				"status": "error",
	//			})()
	//		return
	//	}
	//
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
	clientID := uuid.NewV4().String()

	// Add a listener for popped messages directed to us.
	messageCh := make(chan []byte)
	lid := cc.ListenerAPI.Register(popListener(clientID, messageCh))
	defer close(messageCh)
	defer cc.ListenerAPI.Deregister(lid)

	// Now that we are listening, ask the engine to pop an item off of the
	// queue.
	cc.EventsCh <- engine.Pop{
		ClientID: clientID,
		Channel:  name,
	}

	resp := responder.New(w)

	select {
	case message := <-messageCh:
		resp.Json(map[string]string{
			"status": "ok",
			"data":   string(message),
		})()

	case <-time.After(5 * time.Second):
		resp.Status(http.StatusServiceUnavailable).
			Json(map[string]string{
				"status": "error",
				"error":  "timeout",
			})()
	}
}

func popListener(clientID string, ch chan<- []byte) (engine.EventMatcher, engine.Listener) {
	matcher := engine.EventMatcher{
		EventType: engine.EventPop,
		ClientID:  clientID,
	}

	listener := func(_ engine.Event, changes engine.ChangeSet) {
		ch <- changes.Popped
	}

	return matcher, listener
}
