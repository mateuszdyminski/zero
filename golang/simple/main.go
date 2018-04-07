package main

import (
	"flag"
	"log"
	"net/http"

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

func main() {
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

	log.Fatal(http.ListenAndServe(":"+*port, golang.NewLogginHandler(router)))
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
