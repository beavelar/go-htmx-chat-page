package main

import (
	"log"
)

func main() {
	if err := InitGrpcClient(); err != nil {
		log.Fatalf("error occurred setting up content service grpc client: %s\n", err)
	}
	defer CloseConn()

	if err := InitHttpServer(); err != nil {
		log.Fatalf("error occurred setting up content service http server: %s\n", err)
	}

	log.Println("shutting down content service")
}
