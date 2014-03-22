package main

import (
	"errors"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
)

const VERSION = "1.0-SNAPSHOT"

var configFile = flag.String("c", "", "the location of the server configuration")
var debug *bool = flag.Bool("d", false, "enable debug logging")

// Load the configuration and do the necessary error handling.
func loadConfigFile(path string) (*DTopConfiguration, error) {
	if path == "" {
		return nil, errors.New("Please supply a valid configuration file (-c).")
	}

	if _, err := os.Stat(path); os.IsNotExist(err) {
		return nil, errors.New(fmt.Sprintf("The configuration file does not exist: %s", path))
	}

	jsonBlob, err := ioutil.ReadFile(path)

	if err != nil {
		return nil, errors.New(fmt.Sprintf("Error reading configuration file: %s", err.Error()))
	}

	cfg, err := DeserializeDTopConfigurationFromJson(jsonBlob)

	if err != nil {
		return nil, errors.New(fmt.Sprintf("Error parsing configuration file: %s", err.Error()))
	}

	if valid, err := cfg.IsValid(); !valid {
		return cfg, errors.New(fmt.Sprintf("Invalid configuration: %s", err.Error()))
	}

	return cfg, nil
}

// entry point
func main() {
	// parse cli args
	flag.Parse()

	log.Printf("Reading configuration from '%s'", *configFile)
	cfg, cfgError := loadConfigFile(*configFile)

	if cfgError != nil {
		panic(cfgError)
	}

	auth := NewAuthenticator(cfg)

	// registered publishers
	eventPublishers := [...]EventPublisher{memory, uptime, loadavg, cpuinfo, users, processinfo, basicinfo, disk, services(cfg.Services)}

	// start publishers.
	events := make(chan Event)
	for _, eventPublisher := range eventPublishers {
		go FailSafe(eventPublisher)(events)
	}

	eventServer := NewEventServer(events, jsonEventSerializer)

	// start fanout and monitor goroutines.
	go eventServer.fanOut()
	go eventServer.monitor()

	// register http resources
	http.Handle("/", http.FileServer(http.Dir(cfg.StaticFolder)))
	http.Handle("/events", http.HandlerFunc(eventServer.handler))
	http.Handle("/auth", http.HandlerFunc(authHandler(eventServer, cfg, auth)))

	log.Printf("starting server at 0.0.0.0:%d", cfg.Port)
	err := http.ListenAndServe(fmt.Sprintf(":%d", cfg.Port), nil)

	if err != nil {
		panic("error running server: " + err.Error())
	}
}

// the authHandler function decorator checks for credentials.
func authHandler(eventServer *EventServer, cfg *DTopConfiguration, auth *Authenticator) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method == "POST" {
			username := r.FormValue("username")
			password := r.FormValue("password")

			if success, token := auth.Login(username, password); success {
				eventServer.eventListeners[token] = nil
				w.Write([]byte(token))
			} else {
				log.Printf("received wrong login attempt (user=%s)", username)
				http.Error(w, "bad credentials", 401)
				return
			}
		} else if r.Method == "GET" {
			// TODO: serialize object iso manually creating string here...
			msg := "{\"Name\":\"%s\",\"Auth\":%t,\"Description\":\"%s\",\"Version\":\"%s\"}"
			auth := cfg.IsAuth()
			msg = fmt.Sprintf(msg, cfg.Name, auth, cfg.Description, VERSION)
			w.Write([]byte(msg))
		}
	}
}

// Debug when configured.
func Debugf(format string, args ...interface{}) {
	if *debug {
		log.Printf("DEBUG "+format, args)
	}
}
