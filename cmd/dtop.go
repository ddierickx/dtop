package main

import (
	"flag"
	"fmt"
	"github.com/ddierickx/dtop/pkg"
	"log"
	"net/http"
	"time"
)

const VERSION = "0.1"

var configFile = flag.String("c", "", "the location of the server configuration")
var debug *bool = flag.Bool("d", false, "enable debug logging")

// entry point
func main() {
	// parse cli args
	flag.Parse()

	log.Printf("Reading configuration from '%s'", *configFile)
	cfg, cfgError := pkg.LoadConfigFile(*configFile)

	if cfgError != nil {
		panic(cfgError)
	}

	auth := pkg.NewAuthenticator(cfg)

	// registered publishers
	eventPublishers := [...]pkg.EventPublisher{
		pkg.Repeat(pkg.GetMemory, 2*time.Second),
		pkg.Repeat(pkg.GetUptime, 1*time.Second),
		pkg.Repeat(pkg.GetLoadAvg, 2*time.Second),
		pkg.Repeat(pkg.GetCPUInfo(), 1*time.Second),
		pkg.Repeat(pkg.GetUsers, 3*time.Second),
		pkg.Repeat(pkg.GetProcessInfo, 1*time.Second),
		pkg.Repeat(pkg.GetBasicInfo, 1*time.Hour),
		pkg.Repeat(pkg.GetDisk, 3*time.Second),
		pkg.Repeat(pkg.GetServices(cfg.Services), 2*time.Second),
	}

	// start publishers.
	events := make(chan pkg.Event)
	for _, eventPublisher := range eventPublishers {
		go eventPublisher(events)
	}

	eventServer := pkg.NewEventServer(events, pkg.JsonEventSerializer)

	// start fanout and monitor goroutines.
	go eventServer.FanOut()
	go eventServer.Monitor()

	// register http resources
	http.Handle("/", http.FileServer(http.Dir(cfg.StaticFolder)))
	http.Handle("/events", http.HandlerFunc(eventServer.Handler))
	http.Handle("/auth", http.HandlerFunc(authHandler(eventServer, cfg, auth)))

	log.Printf("starting server at 0.0.0.0:%d", cfg.Port)
	err := http.ListenAndServe(fmt.Sprintf(":%d", cfg.Port), nil)

	if err != nil {
		panic("error running server: " + err.Error())
	}
}

// the authHandler function decorator checks for credentials.
func authHandler(eventServer *pkg.EventServer, cfg *pkg.DTopConfiguration, auth *pkg.Authenticator) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method == "POST" {
			username := r.FormValue("username")
			password := r.FormValue("password")

			if success, token := auth.Login(username, password); success {
				eventServer.EventListeners[token] = nil
				_, _ = w.Write([]byte(token))
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
			_, _ = w.Write([]byte(msg))
		}
	}
}
