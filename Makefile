.DEFAULT_GOAL := help
.PHONY: help lint format generate
SHELL := /bin/bash

.PHONY: help publish local lint test show-acr run
ACR_NAME	?= crmhv25ua3ddectf7dg
KO_DOCKER_REPO	?= $(ACR_NAME).azurecr.io

help:	## show this help
	@grep -E '^[a-zA-Z0-9_-]+:.*?## .*$$' $(MAKEFILE_LIST) | awk 'BEGIN {FS = ":.*?## "}; {printf "%-10s %s\n", $$1, $$2}'

lint:	## run buf lint
lint:	lint-golangci lint-buf

format:	## run buf/Go format
	buf format -w
	go fmt ./...

generate:	## run buf generate
	rm -rf pkg/gen/buf/*
	buf generate

## internal targets

lint-golangci:	# run golangci-lint
	docker run --rm -v $$(pwd):/app -v ~/.cache/golangci-lint/v1.58.0:/root/.cache -w /app golangci/golangci-lint:v1.58.0 golangci-lint run

lint-buf:	# run buf lint
	buf lint



## for local development


local-otelcol-up:  ## up local opentelemetry collector
	docker run \
		-p 127.0.0.1:4317:4317 \
		-p 127.0.0.1:55679:55679 \
		$${OTELCOL_IMG} \
		2>&1 | tee logs/collector-output.txt

local-client:	## local run the client
	go run ./cmd/client http://localhost:8080

local-server:	## local run the server
	go run ./cmd/server

local-grpcurl:	## local run the grpc client
	grpcurl \
	-protoset <(buf build -o -) -plaintext \
	-d '{"message": "hello grpc!"}' \
	localhost:8080 echo.v1.EchoService/Echo

