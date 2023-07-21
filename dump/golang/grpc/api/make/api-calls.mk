.PHONY: client, subscribe, publish

client:
	@go run cmd/broker_client/main.go \
		-grpc-port $(GRPC_PORT)

publish: _http-publish _grpc-publish

_http-publish:
	@echo "HTTP Publish: "
	@curl localhost:$(HTTP_PORT)/mq.Broker/Publish --silent \
  		-X POST -H "Content-Type: application/json" | jq

_grpc-publish:
	@echo "gRPC Publish: "
	@grpcurl -plaintext localhost:$(GRPC_PORT) mq.Broker/Publish | jq

subscribe:
	@grpcurl -d '{"topic": "topic1"}' \
	-plaintext localhost:$(GRPC_PORT) mq.Broker/Subscribe
