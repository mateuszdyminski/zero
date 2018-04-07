package main

import (
	"context"
	"flag"
	"log"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/gorilla/mux"
	"github.com/mateuszdyminski/zero/golang"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

// Flags
var port = flag.String("port", "8080", "HTTP port number")

// Variables injected by -X flag
var appVersion = "unknown"
var gitVersion = "unknown"
var lastCommitTime = "unknown"
var lastCommitHash = "unknown"
var lastCommitUser = "unknown"
var buildTime = "unknown"

// Timeout is the duration to allow outstanding requests to survive before forcefully terminating them.
const Timeout = 10

func main() {
	flag.Parse()

	// subscribe to SIGINT signals
	quit := make(chan os.Signal)
	signal.Notify(quit, os.Interrupt, os.Kill)

	// configure http handlers
	flag.Parse()

	router := mux.NewRouter()

	rest, err := golang.NewUserRest(buildInfo())
	if err != nil {
		log.Fatal(err)
	}

	router.HandleFunc("/api/users", rest.Users).Methods("GET")
	router.HandleFunc("/api/users", rest.AddUser).Methods("POST")
	router.HandleFunc("/api/users/{id}", rest.GetUser).Methods("GET")
	router.HandleFunc("/api/health", rest.Health).Methods("GET")
	router.HandleFunc("/api/error", rest.Err).Methods("POST")
	router.Handle("/metrics", promhttp.HandlerFor(
		prometheus.DefaultGatherer,
		promhttp.HandlerOpts{},
	))

	// create and start http server in new goroutine
	srv := &http.Server{Addr: ":" + *port, Handler: golang.NewLogginHandler(router)}
	go func() {
		// we can't use log.Fatal here!
		if err := srv.ListenAndServe(); err != nil {
			log.Printf("http server stoped: %s\n", err)
		}
	}()

	// blocks the execution until os.Interrupt or os.Kill signal appears
	<-quit

	log.Println("shutting down server. waiting to drain the ongoing requests...")

	// shut down gracefully, but wait no longer than the Timeout value.
	ctx, cancelF := context.WithTimeout(context.Background(), Timeout*time.Second)
	defer cancelF()

	// shutdown the http server
	if err := srv.Shutdown(ctx); err != nil {
		log.Fatalf("error while shutdown http server: %v\n", err)
	}

	log.Println("server gracefully stopped")
}

func buildInfo() golang.BuildInfo {
	return golang.BuildInfo{
		Version:    appVersion,
		GitVersion: gitVersion,
		BuildTime:  buildTime,
		LastCommit: golang.Commit{
			Author: lastCommitUser,
			Id:     lastCommitHash,
			Time:   lastCommitTime,
		},
	}
}
