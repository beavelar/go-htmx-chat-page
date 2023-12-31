package main

import (
	"log"
)

func main() {
	if err := InitDb(); err != nil {
		log.Fatalf("failed to create database connection, exiting - %s\n", err)
	}
	defer CloseDb()

	if err := InitGrpcServer(); err != nil {
		log.Fatalf("failed to start database grpc server, exiting - %s\n", err)
	}

	log.Println("shutting down database service")
}
