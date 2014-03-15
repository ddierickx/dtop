package main

import (
	"os"
	"flag"
	"fmt"
	"log"
	"net/http"
	"github.com/nu7hatch/gouuid"
)

var port = flag.Int("port", 12345, "the tcp port on which to expose the webinterface")
var debug *bool = flag.Bool("debug", false, "enable debug logging")

// entry point
func main() {
	// parse cli args
	flag.Parse()

	// registered publishers
	eventPublishers := [...]EventPublisher{ memory, uptime, loadavg, cpuinfo, users, processinfo, basicinfo, disk }

	// start publishers.
	events := make(chan Event)
	for _, eventPublisher := range eventPublishers {
		go FailSafe(eventPublisher)(events)
	}

	eventServer := NewEventServer(events, jsonEventSerializer)

	// start fanout and monitor goroutines.
	go eventServer.fanOut()
	go eventServer.monitor()

	// register http handlers
	path := "/usr/local/share/dtop/static"
	if _, err := os.Stat("./static"); err == nil {
		path = "./static"
	} else {
		if _, err := os.Stat(path); err != nil {
			panic("Web interface source files were not found at: " + path)
		}
	}

	// register http resources
	http.Handle("/", http.FileServer(http.Dir(path)))
	http.Handle("/events", http.HandlerFunc(eventServer.handler))
	http.Handle("/auth", http.HandlerFunc(authHandler(eventServer)))

	log.Printf("starting server at 0.0.0.0, port %d", *port)
	err := http.ListenAndServe(fmt.Sprintf(":%d", *port), nil)
	
	if err != nil {
		panic("error running server: " + err.Error())
	}
}

// the authHandler function decorator checks for credentials.
func authHandler(eventServer *EventServer) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		username := r.FormValue("username")
		password := r.FormValue("password")

		if username != "" && password != "" {
			Debugf("user '%s' requests authentication", username)

			if username == "u" && password == "p" {
				token, _ := uuid.NewV4()
				eventServer.eventListeners[token.String()] = nil
				w.Write([]byte(token.String()))
			} else {
				log.Printf("received wrong login attempt (user=%s)", username)
				http.Error(w, "bad credentials", 401)
			}
		} else {
			w.Write([]byte(fmt.Sprintf("{\"Server\":\"%s\",\"Auth\":true,\"Description\":\"%s\",\"Version\":\"%s\"}",
								"dominiek-laptop", "Work laptop", "1.0")))
		}
	}
}

// Debug when configured.
func Debugf(format string, args ...interface{}) {
	if *debug {
		log.Printf("DEBUG "+format, args)
	}
}