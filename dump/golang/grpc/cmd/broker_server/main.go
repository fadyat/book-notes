package main

import (
	"context"
	"github.com/fadyat/grpc-broker/api/pb"
	"github.com/fadyat/grpc-broker/internal/broker"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
	"log"
	"net"
)

func main() {
	cfg := parseConfig()
	listener, err := net.Listen("tcp", cfg.GrpcPort())
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	s := grpc.NewServer()
	pb.RegisterBrokerServer(s, broker.NewGrpcServer())

	// Register reflection service on gRPC server.
	// This is helpful for debugging, like grpcurl.
	reflection.Register(s)

	// Registering the gRPC gateway for handling HTTP/1.1 requests.
	// Launching it in a separate goroutine to avoid blocking the
	// gRPC server from serving requests.
	go func() {
		ctx, cancel := context.WithCancel(context.Background())
		defer cancel()

		log.Printf("starting http server on %s", cfg.HTTPPort())
		if e := broker.RunHTTPServer(ctx, cfg.GrpcPort(), cfg.HTTPPort()); e != nil {
			log.Fatalf("failed to serve: %v", e)
		}
	}()

	log.Printf("starting grpc server on %s", cfg.GrpcPort())
	if e := s.Serve(listener); e != nil {
		log.Fatalf("failed to serve: %v", e)
	}
}
