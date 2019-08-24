package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"time"
	"os"
	"os/signal"
	"syscall"

	"github.com/gorilla/mux"
	"gopkg.in/natefinch/lumberjack.v2"
)

func handler(w http.ResponseWriter, r *http.Request) {
	query := r.URL.Query()
	name := query.Get("name")
	if name == "" {
		name = "Guest"
	}
	log.Printf("Received a request for %s\n", name)
	w.Write([]byte(fmt.Sprintf("Hello, %s\n", name)))
}

func main() {
	// create server and route handlers
	r := mux.NewRouter()
	r.HandleFunc("/", handler)

	srv := &http.Server{
		Handler:     r,
		Addr:        ":8080",
		ReadTimeout: 10 * time.Second,
		WriteTimeout: 10 * time.Second,
	}

	// configurable logging
	LOG_FILE_LOCATION := os.Getenv("LOG_FILE_LOCATION")
	if LOG_FILE_LOCATION != "" {
		log.SetOutput(&lumberjack.Logger{
			Filename: LOG_FILE_LOCATION,
			MaxSize: 500,
			MaxBackups: 3,
			MaxAge: 28,
			Compress: true,
		})
	}

	go func() {
		log.Println("starting server")
		if error := srv.ListenAndServe(); error != nil {
			log.Fatal(error)
		}
	}()

	// graceful shutdown
	waitForShutdown(srv)
}

func waitForShutdown(srv *http.Server) {
	interruptChan := make(chan os.Signal, 1)
	signal.Notify(interruptChan, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)

	// Block until we receive a signal.
	<-interruptChan

	// create a deadline to wait for
	ctx, cancel := context.WithTimeout(context.Background(), time.Second * 10)
	defer cancel()
	srv.Shutdown(ctx)

	log.Println("shutting down")
	os.Exit(0)
}
