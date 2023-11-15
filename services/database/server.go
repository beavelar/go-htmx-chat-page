package main

import (
	"log"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	err := InitDb()
	if err != nil {
		log.Fatalf("failed to create database connection, exiting - %s\n", err)
	}
	defer Db.Close()

	err = InitGrpcServer()
	if err != nil {
		log.Fatalf("failed to start database grpc server, exiting - %s\n", err)
	}

	closeC := make(chan os.Signal, 1)
	signal.Notify(closeC, syscall.SIGINT, syscall.SIGTERM)
	<-closeC

	log.Println("shutting down database service")
}
