ifneq (,$(wildcard ./.env.local))
	include .env.local
	export
endif

ifndef HTTP_PORT
	HTTP_PORT=8080
endif

ifndef GRPC_PORT
	GRPC_PORT=8081
endif
