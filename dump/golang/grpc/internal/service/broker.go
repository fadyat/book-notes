package service

import (
	"context"
	"github.com/fadyat/grpc-broker/api/pb"
	"github.com/fadyat/grpc-broker/internal/repo"
)

type Broker interface {

	// Publish publishes a message to a topic.
	Publish(ctx context.Context, in *pb.PublishRequest) (*pb.PublishResponse, error)

	// Subscribe subscribes to a topic.
	Subscribe(ctx context.Context, stream pb.Broker_SubscribeServer) error
}

type broker struct {
	storage repo.Storage
}

func _() Broker {
	return &broker{}
}

func (b *broker) Publish(ctx context.Context, in *pb.PublishRequest) (*pb.PublishResponse, error) {
	panic("implement me")
}

func (b *broker) Subscribe(ctx context.Context, stream pb.Broker_SubscribeServer) error {
	panic("implement me")
}
