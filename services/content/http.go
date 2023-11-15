package main

import (
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"time"
)

func rootReq(w http.ResponseWriter, r *http.Request) {
	msgs, err := GetMessages(0)
	if err != nil {
		io.WriteString(w, "")
		return
	}

	io.WriteString(w, msgs.String())
}

func InitHttpServer() error {
	log.Println("setting up content service http server..")
	host := os.Getenv("CONTENT_SERVICE_HOST")
	if host == "" {
		return errors.New("no CONTENT_SERVICE_HOST environment variable provided, exiting..")
	}

	port := os.Getenv("CONTENT_SERVICE_PORT")
	if port == "" {
		return errors.New("no CONTENT_SERVICE_PORT environment variable provided, exiting..")
	}

	log.Printf("starting up http server: host - %s, port - %s\n", host, port)

	mux := http.NewServeMux()
	httpServer := &http.Server{
		Addr:              fmt.Sprintf("%s:%s", host, port),
		Handler:           mux,
		ReadTimeout:       5 * time.Second,
		ReadHeaderTimeout: 5 * time.Second,
		WriteTimeout:      5 * time.Second,
		IdleTimeout:       10 * time.Second,
	}

	go func() {
		err := httpServer.ListenAndServe()
		log.Fatalf("error occurred attempting to spin up content service http server - %s\n", err)
	}()

	return nil
}
