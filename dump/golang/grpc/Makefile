.PHONY: proto, run

proto:
	@protoc --proto_path=api/proto \
		--go_out=api/pb \
		--go_opt=paths=source_relative \
		--go-grpc_out=api/pb \
		--go-grpc_opt=paths=source_relative \
		api/proto/*.proto

run:
	@go run cmd/broker_server/main.go

client:
	@go run cmd/broker_client/main.go

publish:
	@grpcurl -plaintext localhost:8080 mq.Broker/Publish | jq

subscribe:
	@grpcurl -d '{"topic": "topic1"}' \
	-plaintext localhost:8080 mq.Broker/Subscribe