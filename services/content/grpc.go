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

func GetMessages(ctx context.Context, limit int32) (*proto.Messages, error) {
	if err := checkClient(); err != nil {
		return nil, err
	}

	res, err := client.GetMessages(ctx, &proto.GetMessagesRequest{InitialLimit: &limit})
	if err != nil {
		log.Printf("failed to get messages from database grpc server: %s\n", err)
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

	conn, err := grpc.Dial(fmt.Sprintf("%s:%s", host, port), grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Printf("failed to connect to database service grpc server - %s\n", err)
		return err
	}

	log.Printf("starting up grpc client - host: %s, port: %s\n", host, port)
	client = proto.NewDatabaseServiceClient(conn)
	return nil
}

func StreamMessages(ctx context.Context) (chan *proto.Message, error) {
	if err := checkClient(); err != nil {
		return nil, err
	}

	stream, err := client.StreamMessages(ctx, &proto.StreamMessagesRequest{})
	if err != nil {
		log.Printf("failed to create stream from database grpc server to content service grpc client: %s\n", err)
		return nil, err
	}

	ch := make(chan *proto.Message)
	go func() {
		for {
			select {
			case <-ctx.Done():
				log.Println("received context done, closing message stream channel")
				close(ch)
				return
			default:
				msg, err := stream.Recv()
				if err == io.EOF {
					log.Println("received end of file on database grpc stream, closing channel")
					close(ch)
					return
				}

				if err != nil {
					log.Printf("error occurred receiving messages from database grpc server: %s\n", err)
					close(ch)
					return
				}

				log.Printf("received message from database grpc server, sending out to clients: %s", msg)
				ch <- msg
			}
		}
	}()
	return ch, nil
}

func PostMessage(ctx context.Context, msg *proto.Message) error {
	if err := checkClient(); err != nil {
		return err
	}

	if _, err := client.PostMessage(ctx, msg); err != nil {
		log.Printf("failed to post message to database grpc server: %s\n", err)
		return err
	}

	return nil
}

func checkClient() error {
	if client == nil {
		log.Println("grpc client not initlialized, operation will not be conducted")
		return errors.New("grpc client not initlialized, operation will not be conducted")
	}
	return nil
}
