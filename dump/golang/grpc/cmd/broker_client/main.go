package main

import (
	"context"
	"github.com/fadyat/grpc-broker/api/pb"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"log"
)

const (
	port = ":8080"
)

func main() {
	conn, err := grpc.Dial(
		port, grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	if err != nil {
		log.Fatalf("failed to dial: %v", err)
	}
	defer func() {
		if e := conn.Close(); e != nil {
			log.Fatalf("failed to close connection: %v", e)
		}
	}()

	client := pb.NewBrokerClient(conn)
	publishResponse, err := client.Publish(context.Background(), &pb.PublishRequest{
		Topic: "topic", Body: []byte("hello"), TTL: uint64(1000),
	})
	if err != nil {
		log.Fatalf("failed to publish: %v", err)
	}

	log.Printf("published message with id: %d", publishResponse.Id)
}
