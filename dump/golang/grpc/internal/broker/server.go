package broker

import (
	"context"
	"github.com/fadyat/grpc-broker/api/pb"
	"github.com/fadyat/grpc-broker/pkg"
)

var (
	topics = []string{"topic1", "topic2", "topic3"}
)

type Server struct {
	pb.UnimplementedBrokerServer
}

func NewServer() *Server {
	return &Server{}
}

func (s *Server) Publish(ctx context.Context, in *pb.PublishRequest) (*pb.PublishResponse, error) {
	return &pb.PublishResponse{Id: 1}, nil
}

func (s *Server) Subscribe(in *pb.SubscribeRequest, stream pb.Broker_SubscribeServer) error {
	if pkg.In(topics, in.Topic) {
		return stream.Send(&pb.MessageResponse{Body: []byte("hello")})
	}

	return pkg.ErrorTopicNotFound
}
