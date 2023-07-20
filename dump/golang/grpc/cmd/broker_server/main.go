package main

import (
	"github.com/fadyat/grpc-broker/api/pb"
	"github.com/fadyat/grpc-broker/internal/broker"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
	"log"
	"net"
)

const (
	port = ":8080"
)

func main() {
	listener, err := net.Listen("tcp", port)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	s := grpc.NewServer()
	pb.RegisterBrokerServer(s, broker.NewServer())

	// Register reflection service on gRPC server.
	// This is helpful for debugging, like grpcurl.
	reflection.Register(s)

	log.Println("starting server on port :8080")
	if e := s.Serve(listener); e != nil {
		log.Fatalf("failed to serve: %v", e)
	}
}
