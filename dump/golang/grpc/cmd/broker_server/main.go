package main

import (
	"context"
	"github.com/fadyat/grpc-broker/api/pb"
	"github.com/fadyat/grpc-broker/internal/broker"
	"github.com/fadyat/grpc-broker/internal/logger"
	"github.com/grpc-ecosystem/go-grpc-middleware/v2/interceptors/logging"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
	"net"
)

func main() {
	cfg := parseConfig()
	log := initLogger()
	logOpts := []logging.Option{
		logging.WithLogOnEvents(logging.StartCall, logging.FinishCall),
	}

	s := grpc.NewServer(
		grpc.ChainUnaryInterceptor(
			logging.UnaryServerInterceptor(logger.ToInterceptorLogger(log), logOpts...),
		),
		grpc.ChainStreamInterceptor(
			logging.StreamServerInterceptor(logger.ToInterceptorLogger(log), logOpts...),
		),
	)
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

	listener, err := net.Listen("tcp", cfg.GrpcPort())
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	log.Printf("starting grpc server on %s", cfg.GrpcPort())
	if e := s.Serve(listener); e != nil {
		log.Fatalf("failed to serve: %v", e)
	}
}
