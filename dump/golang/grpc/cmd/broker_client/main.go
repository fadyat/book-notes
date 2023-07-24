package main

import (
	"context"
	"github.com/fadyat/grpc-broker/api/pb"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"log"
)

const (
	grpcPort = ":8081"
)

func main() {
	conn, err := grpc.Dial(
		grpcPort,
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	if err != nil {
		log.Panicf("failed to dial: %v", err)
	}
	defer func() {
		if e := conn.Close(); e != nil {
			log.Printf("failed to close connection: %v", e)
		}
	}()

	client := pb.NewBrokerClient(conn)
	publishResponse, err := client.Publish(context.Background(), &pb.PublishRequest{
		Topic: "topic", Body: []byte("hello"),
	})
	if err != nil {
		log.Panicf("failed to publish: %v", err)
	}

	log.Printf("published message with id: %d", publishResponse.Id)
}
