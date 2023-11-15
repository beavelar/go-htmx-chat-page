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
	client proto.DatabaseServiceClient
	conn   *grpc.ClientConn
)

func CloseConn() {
	if conn != nil {
		conn.Close()
	}
}

func GetMessages(limit int32) (*proto.Messages, error) {
	req := &proto.GetMessagesRequest{}
	if limit > 0 {
		req.InitialLimit = &limit
	}

	if client == nil {
		msg := "unable to get all messages, content service grpc client has not been initialized"
		log.Println(msg)
		return nil, errors.New(msg)
	}

	res, err := client.GetMessages(context.Background(), req)
	if err != nil {
		log.Printf("failed to get messages from database grpc server - %s\n", err)
		return nil, err
	}

	return res, nil
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
	conn, err = grpc.Dial(fmt.Sprintf("%s:%s", host, port), grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Printf("failed to connect to database service grpc server - %s\n", err)
		return err
	}

	client = proto.NewDatabaseServiceClient(conn)
	return nil
}

func StreamMessages() (chan *proto.Message, error) {
	if client == nil {
		msg := "unable to stream messages, content service grpc client has not been initialized"
		log.Println(msg)
		return nil, errors.New(msg)
	}

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
	if client == nil {
		msg := "unable to post message, content service grpc client has not been initialized"
		log.Println(msg)
		return errors.New(msg)
	}

	_, err := client.PostMessage(context.Background(), msg)
	if err != nil {
		log.Printf("failed to post message to database grpc server - %s\n", err)
		return err
	}

	return nil
}
