package main

import (
	"log"
)

func main() {
	err := InitGrpcClient()
	if err != nil {
		log.Fatalf("error occurred setting up content service grpc client: %s\n", err)
	}
	defer CloseConn()

	err = InitHttpServer()
	if err != nil {
		log.Fatalf("error occurred setting up content service http server: %s\n", err)
	}

	log.Println("shutting down content service")
}
