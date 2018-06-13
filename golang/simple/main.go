package main

import (
	"flag"
	"log"
	"net/http"

	_ "github.com/go-sql-driver/mysql"
	"github.com/gorilla/mux"
	"github.com/mateuszdyminski/zero/golang"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

// Flags
var port = flag.String("port", "8080", "HTTP port number")
var dbHost = flag.String("db-host", "mysql-mysql:3306", "Mysql database host with port")
var dbUser = flag.String("db-user", "root", "Mysql database user")
var dbPass = flag.String("db-pass", "password", "Mysql database password")
var dbName = flag.String("db-name", "users", "Mysql database name")

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

	rest, err := golang.NewUserRest(buildInfo(), dbInfo())
	if err != nil {
		log.Fatal(err)
	}

	router.HandleFunc("/api/users", rest.Users).Methods("GET")
	router.HandleFunc("/api/users", rest.AddUser).Methods("POST")
	router.HandleFunc("/api/users/{id}", rest.GetUser).Methods("GET")
	router.HandleFunc("/health", rest.Health).Methods("GET")
	router.HandleFunc("/api/error", rest.Err).Methods("POST")
	router.Handle("/metrics", promhttp.HandlerFor(
		prometheus.DefaultGatherer,
		promhttp.HandlerOpts{},
	))

	log.Println("starting http server!")

	log.Fatal(http.ListenAndServe(":"+*port, golang.NewLogginHandler(router)))
}

func buildInfo() golang.BuildInfo {
	return golang.BuildInfo{
		Version:    appVersion,
		GitVersion: gitVersion,
		BuildTime:  buildTime,
		LastCommit: golang.Commit{
			Author: lastCommitUser,
			ID:     lastCommitHash,
			Time:   lastCommitTime,
		},
	}
}

func dbInfo() golang.DBInfo {
	return golang.DBInfo{
		Name:     *dbName,
		Host:     *dbHost,
		User:     *dbUser,
		Password: *dbPass,
	}
}
