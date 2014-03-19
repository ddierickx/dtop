package main

import (
	"github.com/gorilla/websocket"
	"log"
	"net/http"
	"sync"
	"time"
)

// The eventserver will read from the publishers and fanout events to all connected clients.
type EventServer struct {
	eventsCount         int64
	events              chan Event // input channel
	eventListenersMutex sync.RWMutex
	eventListeners      map[string]chan Event // fan-out channels
	lastEvents          map[string]Event
	eventSerializer     EventSerializer
}

// The monitor function of the eventserver outputs debug info on the rate of published events.
func (eventServer *EventServer) monitor() {
	for {
		start := eventServer.eventsCount
		time.Sleep(1 * time.Minute)
		delta := eventServer.eventsCount - start
		log.Printf("rate %de/m", delta)
	}
}

// Receive from each channel and fanout to each connected client channel.
func (eventServer *EventServer) fanOut() {
	for event := range eventServer.events {
		// lock client connection channels list.
		eventServer.eventListenersMutex.RLock()
		eventServer.lastEvents[event.Q] = event
		eventServer.eventsCount += 1

		for _, listener := range eventServer.eventListeners {
			if listener != nil { // unregistered connected channels have nil values.
				eventServer.fanOutSafe(listener, event)
			}
		}

		eventServer.eventListenersMutex.RUnlock()
	}
}

// Try-catch around channel submission to handle potential deadlock errors upon disconnect.
func (eventServer *EventServer) fanOutSafe(listener chan Event, event Event) {
	defer func() {
		if err := recover(); err != nil {
			// swallow closed channel error
			// http://stackoverflow.com/questions/16105325/how-to-check-a-channel-is-closed-or-not-without-reading-it
		}
	}()
	listener <- event
}

// Register a new client connection channel.
func (eventServer *EventServer) register(token string) chan Event {
	listener := make(chan Event)
	eventServer.eventListeners[token] = listener
	return listener
}

// Unregister the client connection upon disconnected by setting to nil.
func (eventServer *EventServer) unregister(token string) {
	close(eventServer.eventListeners[token])
	delete(eventServer.eventListeners, token)
}

// EventServer constructor.
func NewEventServer(events chan Event, eventSerializer EventSerializer) *EventServer {
	eventServer := new(EventServer)
	eventServer.events = events
	eventServer.eventSerializer = eventSerializer
	eventServer.eventListeners = make(map[string]chan Event)
	eventServer.lastEvents = make(map[string]Event)
	return eventServer
}

// Function that flushes all of the last events to a listener and signals the
// client that it can show the interface (sig.ready).
func (eventServer *EventServer) submitLastEvents(listener chan Event) {
	for _, event := range eventServer.lastEvents {
		listener <- event
	}
	listener <- NewEvent("sig.ready", "")
}

// Client handler, register, monitor channel and transmit, unregister.
func (eventServer *EventServer) handler(w http.ResponseWriter, r *http.Request) {
	ws, err := websocket.Upgrade(w, r, nil, 1024, 1024)

	if _, ok := err.(websocket.HandshakeError); ok {
		http.Error(w, "not a valid websocket handshake", 400)
		return
	} else if err != nil {
		log.Println(err)
		return
	}

	_, tokenBytes, _ := ws.ReadMessage()
	token := string(tokenBytes)

	if _, ok := eventServer.eventListeners[token]; !ok {
		http.Error(w, "unauthorized access", 201)
		return
	}

	listener := eventServer.register(token)
	log.Printf("client %s succesfully connected (token=%s)", ws.RemoteAddr(), token)

	go eventServer.submitLastEvents(listener)

	for {
		item, open := <-listener

		if !open {
			break
		}

		data := eventServer.eventSerializer(item)
		Debugf("sent event: %s", data)

		if err := ws.WriteMessage(websocket.TextMessage, data); err != nil {
			break
		}
	}
	eventServer.unregister(token)
	log.Printf("client %s disconnected (token=%s)", ws.RemoteAddr(), token)
}
