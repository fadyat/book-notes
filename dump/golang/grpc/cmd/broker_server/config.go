package main

import (
	"flag"
	"strconv"
)

type config struct {
	grpcPort int
	httpPort int
}

func getPort(port int) string {
	return ":" + strconv.Itoa(port)
}

func (c *config) GrpcPort() string {
	return getPort(c.grpcPort)
}

func (c *config) HTTPPort() string {
	return getPort(c.httpPort)
}

func parseConfig() *config {
	grpcPort := flag.Int("grpc-port", 8081, "gRPC port for serving")
	httpPort := flag.Int("http-port", 8080, "HTTP port for serving")

	flag.Parse()
	return &config{
		grpcPort: *grpcPort,
		httpPort: *httpPort,
	}
}
