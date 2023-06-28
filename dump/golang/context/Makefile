server:
	@go run slow_server.go

lint:
	@golangci-lint run *.go

with-timeout:
	@go run with_timeout.go client.go

with-cancel:
	@go run with_cancel.go client.go

with-value:
	@go run with_value.go client.go

.PHONY: server, lint, with-timeout, with-cancel, with-value