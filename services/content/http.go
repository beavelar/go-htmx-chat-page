package main

import (
	proto "content-service/genproto/database"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/gorilla/websocket"
)

var (
	pHost    string
	upgrader = websocket.Upgrader{
		CheckOrigin:      checkOrigin,
		HandshakeTimeout: 5 * time.Second,
	}
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

	pHost = os.Getenv("PROXY_HOST")
	if pHost == "" {
		return errors.New("no PROXY_HOST environment variables provided, exiting..")
	}

	log.Printf("starting up http server - host: %s, port: %s\n", host, port)

	mux := http.NewServeMux()
	mux.HandleFunc("/message", message)
	mux.HandleFunc("/messages", messages)
	mux.HandleFunc("/ws/chat", chat)

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

func chat(w http.ResponseWriter, r *http.Request) {
	log.Println("received messages websocket connection request..")
	conn, err := upgrader.Upgrade(w, r, nil)

	if err != nil {
		log.Printf("failed to upgrade http connection to a websocket connection: %s\n", err)
		return
	}
	defer conn.Close()

	pCh := pingPong(r.Context(), conn)
	sCh, err := StreamMessages(r.Context())
	if err != nil {
		log.Printf("failed to get stream of messages, websocket will be closed: %s\n", err)
		return
	}

	for {
		select {
		case _, ok := <-pCh:
			if !ok {
				log.Println("ping pong channel closed, closing websocket connection")
				return
			}
		case msg, ok := <-sCh:
			if !ok {
				log.Println("stream messages channel closed, closing websocket connection")
				return
			}

			log.Printf("received message from grpc channel: %s\n", msg)
			component := MessageUpdate(msg)
			wr, err := conn.NextWriter(websocket.TextMessage)
			if err != nil {
				log.Printf("failed to get a websocket writer, message update won't be sent: %s\n", err)
				continue
			}

			component.Render(r.Context(), wr)
			if err = wr.Close(); err != nil {
				log.Printf("failed to close message writer, message may be delayed: %s\n", err)
			}
		case <-r.Context().Done():
			log.Println("context done received, closing websocket connection")
			return
		}
	}
}

func checkOrigin(r *http.Request) bool {
	return r.Host == pHost
}

func message(w http.ResponseWriter, r *http.Request) {
	log.Println("received message request..")
	if r.Method != http.MethodPost {
		log.Printf("unsupported http method used, request will be ignored")
		w.Header().Set("Allow", "POST")
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	r.ParseForm()
	if !r.Form.Has("message") {
		log.Printf("form doesn't contain message field, request will be ignored")
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	msg := r.Form.Get("message")
	if strings.Trim(msg, " ") == "" {
		log.Println("empty string provided, message will be ignored")
		component := MessageInput()
		component.Render(r.Context(), w)
		return
	}

	err := PostMessage(r.Context(), &proto.Message{Message: msg, Name: "Unknown User", Time: time.Now().Unix()})
	if err != nil {
		log.Printf("failed to post message: %s\n", err)
		w.WriteHeader(http.StatusInternalServerError)
	}

	component := MessageInput()
	component.Render(r.Context(), w)
}

func messages(w http.ResponseWriter, r *http.Request) {
	log.Println("received messages request..")

	res, err := GetMessages(r.Context(), 0)
	if err != nil {
		log.Printf("error occurred attempting to get all messages: %s\n", err)
		res = &proto.Messages{Messages: make([]*proto.Message, 0)}
	}

	component := Messages(res)
	component.Render(r.Context(), w)
}

func pingPong(ctx context.Context, conn *websocket.Conn) chan int {
	ch := make(chan int)
	go func() {
		for {
			select {
			case <-ctx.Done():
				log.Println("received context done, closing ping pong channel")
				close(ch)
				return
			default:
				mt, wsB, err := conn.ReadMessage()
				if err != nil {
					log.Printf("failed to receive message from websocket connection: %s\n", err)
					close(ch)
					return
				}

				var wsMsg map[string]interface{}
				err = json.Unmarshal(wsB, &wsMsg)
				if err != nil {
					log.Printf("failed to decode websocket message, ignoring message: %s\n", err)
					continue
				}

				if wsMsg["ping"] == "ping" {
					log.Println("received ping message from websocket connection, responding with pong..")
					err = conn.WriteMessage(mt, []byte("pong"))
					if err != nil {
						log.Printf("failed to respond to ping message with pong, closing websocket connection: %s\n", err)
						close(ch)
						return
					}
				}
			}
		}
	}()
	return ch
}
