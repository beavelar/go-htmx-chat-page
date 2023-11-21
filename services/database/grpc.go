package main

import (
	"context"
	proto "database-service/genproto/database"
	"errors"
	"fmt"
	"log"
	"net"
	"os"
	"sync"

	"google.golang.org/grpc"
)

type dbGrpcServer struct {
	proto.UnimplementedDatabaseServiceServer
}

var (
	msgCh chan *proto.Message
	mut   = &sync.Mutex{}
	subs  = make(map[*context.Context]*chan *proto.Message)
)

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
	msgCh = make(chan *proto.Message)

	grpcServer := grpc.NewServer()
	proto.RegisterDatabaseServiceServer(grpcServer, &dbGrpcServer{})
	grpcServer.Serve(lis)

	closeMsgCh()
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
	if err := checkMsgCh(); err != nil {
		return err
	}

	ctx := stream.Context()
	ch := make(chan *proto.Message)
	subs[&ctx] = &ch

	for {
		select {
		case msg, ok := <-ch:
			if !ok {
				log.Println("message channel has been closed, ending grpc stream")
				return errors.New("message channel has been closed, ending grpc stream")
			}

			log.Printf("received new messsage from channel, sending out to clients: %s\n", msg)
			if err := stream.Send(msg); err != nil {
				log.Printf("failed to send message through grpc stream, ending grpc stream: %s\n", err)
				return err
			}
		case <-ctx.Done():
			log.Println("grpc stream has been ended")
			close(ch)
			mut.Lock()
			delete(subs, &ctx)
			mut.Unlock()
			return nil
		}
	}
}

func (s *dbGrpcServer) PostMessage(ctx context.Context, req *proto.Message) (*proto.PostMessageResponse, error) {
	log.Printf("received PostMessage request with the following message - message: %s, name: %s, time: %d\n", req.Message, req.Name, req.Time)
	if err := PostMessage(req); err != nil {
		log.Printf("error occurred posting the following message - message: %s, name: %s, time: %d - %s\n", req.Message, req.Name, req.Time, err)
		return nil, err
	}

	if err := sendMessage(req); err != nil {
		return nil, err
	}

	return &proto.PostMessageResponse{Success: true}, nil
}

func checkMsgCh() error {
	if msgCh == nil {
		log.Println("grpc message channel has not been initialized, operation will not be conducted")
		return errors.New("grpc message channel has not been initialized, operation will not be conducted")
	}
	return nil
}

func closeMsgCh() {
	close(msgCh)
	msgCh = nil
}

func sendMessage(msg *proto.Message) error {
	if err := checkMsgCh(); err != nil {
		return err
	}

	for _, ch := range subs {
		*ch <- msg
	}

	return nil
}
