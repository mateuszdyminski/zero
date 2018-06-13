package main

import (
	"context"
	"flag"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

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

// Timeout is the duration to allow outstanding requests to survive before forcefully terminating them.
const Timeout = 20

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

	log.Println("starting http server with graceful shutdown mode!")

	// create and start http server in new goroutine
	srv := &http.Server{Addr: ":" + *port, Handler: golang.NewLogginHandler(router)}
	go func() {
		// we can't use log.Fatal here!
		if err := srv.ListenAndServe(); err != nil {
			log.Printf("http server stoped: %s\n", err)
		}
	}()

	// subscribe to SIGTERM signals
	quit := make(chan os.Signal)
	signal.Notify(quit, os.Interrupt, syscall.SIGTERM)

	// blocks the execution until os.Interrupt or syscall.SIGTERM signal appears
	<-quit
	log.Println("shutting down server. waiting to drain the ongoing requests...")
	rest.Unhealthy()

	// add extra time to prevent new requests be routed to our service
	time.Sleep(5 * time.Second)

	// shut down gracefully, but wait no longer than the Timeout value.
	ctx, cancelF := context.WithTimeout(context.Background(), Timeout*time.Second)
	defer cancelF()

	// shutdown the http server
	if err := srv.Shutdown(ctx); err != nil {
		log.Printf("error while shutdown http server: %v\n", err)
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
