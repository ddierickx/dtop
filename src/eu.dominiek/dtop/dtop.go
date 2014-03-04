package main

import (
	"code.google.com/p/go.net/websocket"
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"net/http"
	"sync"
	"time"
)

var port = flag.Int("port", 12345, "the tcp port on which to expose the webinterface")
var debug *bool = flag.Bool("debug", false, "enable debug logging")

// entry point
func main() {
	// parse cli args
	flag.Parse()

	events := make(chan Event)

	// registered publishers
	eventPublishers := [...]EventPublisher{memory, uptime, loadavg, cpuinfo, users, processinfo, basicinfo, disk}

	// start publishers as parallel.
	for _, eventPublisher := range eventPublishers {
		go FailSafe(eventPublisher)(events)
	}

	eventServer := NewEventServer(events, jsonEventSerializer)

	// start fanout and monitor goroutines.
	go eventServer.fanOut()
	go eventServer.monitor()

	// register http handlers
	http.Handle("/", http.FileServer(http.Dir("./static")))
	http.Handle("/events", websocket.Handler(eventServer.handler))

	log.Printf("starting server at http://127.0.0.1:%d", *port)
	err := http.ListenAndServe(fmt.Sprintf(":%d", *port), nil)
	if err != nil {
		panic("ListenAndServe: " + err.Error())
	}
}

func Debugf(format string, args ...interface{}) {
	if *debug {
		log.Printf("DEBUG "+format, args)
	}
}

// The eventserver will read from the publishers and fanout events to all connected clients.
type EventServer struct {
	eventsCount         int64
	events              chan Event // input channel
	eventListenersMutex sync.RWMutex
	eventListeners      []chan Event // fan-out channels
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
				fanOutSafe(listener, event)
			}
		}

		eventServer.eventListenersMutex.RUnlock()
	}
}

func fanOutSafe(listener chan Event, event Event) {
	defer func() {
		if err := recover(); err != nil {
			// swallow closed channel error
			// http://stackoverflow.com/questions/16105325/how-to-check-a-channel-is-closed-or-not-without-reading-it
		}
	}()
	listener <- event
}

// register a new client connection channel.
func (eventServer *EventServer) register() (chan Event, int) {
	listener := make(chan Event)
	eventServer.eventListenersMutex.Lock()
	eventServer.eventListeners = append(eventServer.eventListeners, listener)
	listenerId := len(eventServer.eventListeners) - 1
	eventServer.eventListenersMutex.Unlock()
	return listener, listenerId
}

// unregister the client connection upon disconnected by setting to nil.
func (eventServer *EventServer) unregister(listenerId int) {
	// close before aquiring lock so pushes from fanOut don't deadlock the channel.
	close(eventServer.eventListeners[listenerId])
	eventServer.eventListenersMutex.Lock()
	eventServer.eventListeners[listenerId] = nil
	eventServer.eventListenersMutex.Unlock()
}

// Eventserver constructor.
func NewEventServer(events chan Event, eventSerializer EventSerializer) *EventServer {
	eventServer := new(EventServer)
	eventServer.events = events
	eventServer.eventSerializer = eventSerializer
	eventServer.lastEvents = make(map[string]Event)
	return eventServer
}

func (eventServer *EventServer) submitLastEvents(listener chan Event) {
	for _, event := range eventServer.lastEvents {
		listener <- event
	}
	listener <- NewEvent("sig.ready", "")
}

// Client handler, register, monitor channel and transmit, unregister.
func (eventServer *EventServer) handler(ws *websocket.Conn) {
	listener, listenerId := eventServer.register()
	// TODO: get actual client ip address.
	log.Printf("client connect %s (id=%d)", ws.Request().RemoteAddr, listenerId)

	go eventServer.submitLastEvents(listener)

	for {
		item, open := <-listener

		if !open {
			break
		}

		data := eventServer.eventSerializer(item)
		Debugf("sent event: %s", data)

		if _, err := ws.Write(data); err != nil {
			break
		}
	}
	eventServer.unregister(listenerId)
	log.Printf("client disconnect %s (id=%d)", ws.Request().RemoteAddr, listenerId)
}

// Server->client serialization protocol function template.
type EventSerializer func(Event) []byte

// Default JSON implementation of EventSerializer.
func jsonEventSerializer(event Event) []byte {
	serialized, err := json.Marshal(event)

	if err != nil {
		panic("Serialization error: " + err.Error())
	}

	return serialized
}
