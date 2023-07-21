include api/make/api-calls.mk
include api/make/env.mk
include api/make/help.mk

pre: proto lint test ##@git Run pre-commit checks.

proto: ##@api Generate api based on proto files.
	@echo "Generating api based on proto files..."
	@protoc --proto_path api/proto \
		--go_out api/pb \
		--go_opt paths=source_relative \
		--go-grpc_out api/pb \
		--go-grpc_opt paths=source_relative \
		--grpc-gateway_out api/pb \
		--grpc-gateway_opt paths=source_relative \
		--grpc-gateway_opt logtostderr=true \
		--grpc-gateway_opt generate_unbound_methods=true \
		--openapiv2_out ./api/openapi \
    	--openapiv2_opt logtostderr=true \
    	--openapiv2_opt generate_unbound_methods=true \
		api/proto/*.proto
	@echo "Done."

lint: ##@api Run linter.
	@echo "Running linter..."
	@golangci-lint run ./...
	@echo "Done."

test: ##@api Run tests.
	@echo "Running tests..."
	@go test --cover ./cmd/... ./internal/... ./pkg/...
	@echo "Done."

run: ##@api Run broker gRPC and HTTP servers.
	@go run cmd/broker_server/*.go \
		--http-port $(HTTP_PORT) \
		--grpc-port $(GRPC_PORT)


.PHONY: pre, proto, lint, test, run
