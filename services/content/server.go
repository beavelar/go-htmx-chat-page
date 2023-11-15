package main

import (
	"log"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	err := InitGrpcClient()
	if err != nil {
		log.Fatalf("error occurred setting up content service grpc client - %s\n", err)
	}

	err = InitHttpServer()
	if err != nil {
		log.Fatalf("error occurred setting up content service http server - %s\n", err)
	}

	closeC := make(chan os.Signal, 1)
	signal.Notify(closeC, syscall.SIGINT, syscall.SIGTERM)
	<-closeC

	log.Println("shutting down content service")
	defer GrpcConn.Close()
}
