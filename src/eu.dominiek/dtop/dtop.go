package main

import (
	"os"
	"flag"
	"fmt"
	"log"
	"io/ioutil"
	"github.com/nu7hatch/gouuid"
	"net/http"
)

var configFile = flag.String("c", "", "the location of the server configuration")
var debug *bool = flag.Bool("d", false, "enable debug logging")

// entry point
func main() {
	// parse cli args
	flag.Parse()

	if *configFile == "" {
		panic("Please supply a valid configuration file (-c).")
	}

	if _, err := os.Stat(*configFile); os.IsNotExist(err) {
    	panic(fmt.Sprintf("The configuration file does not exist: %s", configFile))
	}

	log.Printf("Reading configuration from '%s'", *configFile)
	jsonBlob, err := ioutil.ReadFile(*configFile)
    
    if err != nil {
    	panic(fmt.Sprintf("Error reading configuration file: %s", err.Error()))
    }

	cfg, err := DeserializeDTopConfigurationFromJson(jsonBlob)

	if err != nil {
		panic(fmt.Sprintf("Error parsing configuration file: %s", err.Error()))
	}

	if valid, err := cfg.IsValid(); !valid {
		panic(fmt.Sprintf("Invalid configuration: %s", err.Error()))
	}

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
	http.Handle("/", http.FileServer(http.Dir(cfg.StaticFolder)))
	http.Handle("/events", http.HandlerFunc(eventServer.handler))
	http.Handle("/auth", http.HandlerFunc(authHandler(eventServer)))

	log.Printf("starting server at 0.0.0.0, port %d", cfg.Port)
	err = http.ListenAndServe(fmt.Sprintf(":%d", cfg.Port), nil)
	
	if err != nil {
		panic("error running server: " + err.Error())
	}
}

// the authHandler function decorator checks for credentials.
func authHandler(eventServer *EventServer) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method == "POST" {
			username := r.FormValue("username")
			password := r.FormValue("password")
			if username != "" && password != "" {
				log.Printf("user '%s' requests authentication", username)
				token, _ := uuid.NewV4()
				if username == "u" && password == "p" {
					eventServer.eventListeners[token.String()] = nil
					w.Write([]byte(token.String()))
				} else {
					log.Printf("received wrong login attempt (user=%s)", username)
					http.Error(w, "bad credentials", 401)
					return
				}
			} else {
				http.Error(w, "bad credentials", 401)
				return
			}
		} else if r.Method == "GET" {
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