package broker

import (
	"context"
	"github.com/fadyat/grpc-broker/api/pb"
	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"net/http"
	"time"
)

func RunHTTPServer(ctx context.Context, grpcPort, httpPort string) error {
	mux := runtime.NewServeMux()
	opts := []grpc.DialOption{grpc.WithTransportCredentials(insecure.NewCredentials())}

	if err := pb.RegisterBrokerHandlerFromEndpoint(ctx, mux, grpcPort, opts); err != nil {
		return err
	}

	return runHTTPServer(httpPort, mux)
}

func runHTTPServer(httpPort string, mux *runtime.ServeMux) error {
	server := &http.Server{
		Handler:      mux,
		Addr:         httpPort,
		ReadTimeout:  3 * time.Second,
		WriteTimeout: 3 * time.Second,
	}

	return server.ListenAndServe()
}
