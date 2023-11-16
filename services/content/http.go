package main

import (
	proto "content-service/genproto/database"
	"errors"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"
)

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

	log.Printf("starting up http server - host: %s, port: %s\n", host, port)

	mux := http.NewServeMux()
	mux.HandleFunc("/message", message)
	mux.HandleFunc("/messages", messages)

	httpServer := &http.Server{
		Addr:              fmt.Sprintf("%s:%s", host, port),
		Handler:           mux,
		ReadTimeout:       5 * time.Second,
		ReadHeaderTimeout: 5 * time.Second,
		WriteTimeout:      5 * time.Second,
		IdleTimeout:       10 * time.Second,
	}

	err := httpServer.ListenAndServe()
	log.Fatalf("error occurred attempting to spin up content service http server: %s\n", err)

	return nil
}

func message(w http.ResponseWriter, r *http.Request) {
	log.Println("received message request..")
	if r.Method == http.MethodPost {
		r.ParseForm()

		if r.Form.Has("message") {
			err := PostMessage(&proto.Message{Message: r.Form.Get("message"), Name: "UnkownUser", Time: time.Now().Unix()})
			if err != nil {
				log.Printf("failed to post message: %s\n", err)
				w.WriteHeader(http.StatusInternalServerError)
			}
		}
	}

	component := Message()
	component.Render(r.Context(), w)
}

func messages(w http.ResponseWriter, r *http.Request) {
	log.Println("received messages request..")

	res, err := GetMessages(0)
	if err != nil {
		log.Printf("error occurred attempting to get all messages: %s\n", err)
		res = &proto.Messages{Messages: make([]*proto.Message, 0)}
	}

	component := Messages(res)
	component.Render(r.Context(), w)
}
