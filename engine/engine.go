package engine

import (
	"errors"
)

var (
	ErrNoRoute     = errors.New("Route not found.")
	ErrBufferFull  = errors.New("Route buffer full.")
	ErrBufferEmpty = errors.New("Route buffer empty.")
)

type Engine struct {
	// maxBuf is the number of messages to buffer before dropping
	// messages.
	maxBuf int

	routes map[string]*Route
}

func New(maxBuf int) *Engine {
	return &Engine{
		maxBuf: maxBuf,
		routes: make(map[string]*Route),
	}
}

func (eng *Engine) AddChannel(channelName string) {
	eng.routes[channelName] = &Route{
		ch: make(chan []byte, eng.maxBuf),
	}
}

func (eng *Engine) RemoveChannel(channelName string) {
	route, ok := eng.routes[channelName]
	if ok {
		route.Close()
		eng.routes[channelName] = nil
	}
}

func (eng *Engine) Publish(channelName string, message []byte) error {
	route, ok := eng.routes[channelName]
	if !ok {
		return ErrNoRoute
	}

	select {
	case route.ch <- message:
		return nil
	default:
		return ErrBufferFull
	}
}

func (eng *Engine) Listen(channelName string) <-chan []byte {
	route, ok := eng.routes[channelName]
	if !ok {
		// Make and return a closed channel if the requested channel does not
		// exist in the engine.
		ch := make(chan []byte)
		close(ch)
		return ch
	}

	return route.ch
}

func (eng *Engine) Pop(channelName string) ([]byte, error) {
	route, ok := eng.routes[channelName]
	if !ok {
		return nil, ErrNoRoute
	}

	select {
	case message := <-route.ch:
		return message, nil
	default:
		return nil, ErrBufferEmpty
	}
}
