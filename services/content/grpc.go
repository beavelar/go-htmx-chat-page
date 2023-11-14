package main

import (
	proto "content-service/genproto/database"
	"context"
	"errors"
	"fmt"
	"io"
	"log"
	"os"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

var (
	client   proto.DatabaseServiceClient
	GrpcConn *grpc.ClientConn
)

func GetMessages(limit int32) (*proto.Messages, error) {
	req := &proto.GetMessagesRequest{}
	if limit > 0 {
		req.InitialLimit = &limit
	}

	res, err := client.GetMessages(context.Background(), req)
	if err != nil {
		log.Printf("failed to get messages from database grpc server - %s\n", err)
		return nil, err
	}

	return res, nil
}

func StreamMessages() (chan *proto.Message, error) {
	stream, err := client.StreamMessages(context.Background(), &proto.StreamMessagesRequest{})
	if err != nil {
		log.Printf("failed to create stream from database grpc server to content service grpc client - %s\n", err)
		return nil, err
	}

	c := make(chan *proto.Message)

	go func() {
		for {
			msg, err := stream.Recv()
			if err == io.EOF {
				break
			}

			if err != nil {
				log.Printf("error occurred receiving messages from database grpc server - %v\n", err)
				return
			}

			c <- msg
		}
	}()

	return c, nil
}

func PostMessage(msg *proto.Message) error {
	_, err := client.PostMessage(context.Background(), msg)
	if err != nil {
		log.Printf("failed to post message to database grpc server - %s\n", err)
		return err
	}

	return nil
}

func InitGrpcClient() error {
	log.Println("setting up content service grpc client..")
	host := os.Getenv("DB_SERVICE_HOST")
	if host == "" {
		return errors.New("no DB_SERVICE_HOST environment variable provided, exiting..")
	}

	port := os.Getenv("DB_SERVICE_PORT")
	if port == "" {
		return errors.New("no DB_SERVICE_PORT environment variable provided,  exiting..")
	}

	var err error
	GrpcConn, err = grpc.Dial(fmt.Sprintf("%s:%s", host, port), grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Printf("failed to connect to database service grpc server - %s\n", err)
		return err
	}

	client = proto.NewDatabaseServiceClient(GrpcConn)
	return nil
}
