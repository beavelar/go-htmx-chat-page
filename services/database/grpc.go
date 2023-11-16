package main

import (
	"context"
	proto "database-service/genproto/database"
	"errors"
	"fmt"
	"log"
	"net"
	"os"
	"time"

	"google.golang.org/grpc"
)

type dbGrpcServer struct {
	proto.UnimplementedDatabaseServiceServer
}

func InitGrpcServer() error {
	log.Println("setting up database grpc server..")
	host := os.Getenv("DB_GRPC_HOST")
	if host == "" {
		return errors.New("no DB_GRPC_HOST environment variable provided, exiting..")
	}

	port := os.Getenv("DB_GRPC_PORT")
	if port == "" {
		return errors.New("no DB_GRPC_PORT environment variable provided,  exiting..")
	}

	lis, err := net.Listen("tcp", fmt.Sprintf("%s:%s", host, port))
	if err != nil {
		log.Fatalf("failed to start database grpc server on: %s:%s - %s", host, port, err)
	}

	log.Printf("starting up database grpc server on %s:%s\n", host, port)

	grpcServer := grpc.NewServer()
	proto.RegisterDatabaseServiceServer(grpcServer, &dbGrpcServer{})
	grpcServer.Serve(lis)

	return nil
}

func (s *dbGrpcServer) GetMessages(ctx context.Context, req *proto.GetMessagesRequest) (*proto.Messages, error) {
	log.Printf("received GetMessages request with limit: %d\n", *req.InitialLimit)
	msgs, err := GetMessages(*req.InitialLimit)
	if err != nil {
		log.Printf("error occurred retrieving all messages with limit: %d - %s\n", *req.InitialLimit, err)
		return nil, err
	}
	return msgs, nil
}

func (s *dbGrpcServer) StreamMessages(req *proto.StreamMessagesRequest, stream proto.DatabaseService_StreamMessagesServer) error {
	log.Println("received StreamMessages request")
	for {
		time.Sleep(time.Second)
		if err := stream.Send(&proto.Message{}); err != nil {
			return err
		}
	}
}

func (s *dbGrpcServer) PostMessage(ctx context.Context, req *proto.Message) (*proto.PostMessageResponse, error) {
	log.Printf("received PostMessage request - message: %s\n", req)
	err := PostMessage(req)
	if err != nil {
		log.Printf("error occurred posting the following message - message: %s, name: %s, time: %d - %s\n", req.Message, req.Name, req.Time, err)
		return nil, err
	}
	return &proto.PostMessageResponse{Success: true}, nil
}
