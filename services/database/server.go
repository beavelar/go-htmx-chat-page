package main

import (
	"log"
)

func main() {
	err := InitDb()
	if err != nil {
		log.Fatalf("failed to create database connection, exiting - %s\n", err)
	}
	defer CloseDb()

	err = InitGrpcServer()
	if err != nil {
		log.Fatalf("failed to start database grpc server, exiting - %s\n", err)
	}

	log.Println("shutting down database service")
}
