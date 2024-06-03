.DEFAULT_GOAL := help
.PHONY: help lint format generate
SHELL := /bin/bash

.PHONY: help publish local lint test show-acr run
KO_DOCKER_REPO	?= $(ACR_NAME).azurecr.io

-include .env

## git tagではなく、devでのビルド番号を使う

TAG_VERSION		?= 0.0.1-dev-$(shell date +%Y%m%d)
BUILD_NUMBER_FILE	?= build_number_$(TAG_VERSION).txt
TAG_COUNT		?= $(shell cat $(BUILD_NUMBER_FILE))
IMAGE_NAME		?= $(KO_DOCKER_REPO)/server:$(TAG_VERSION)-$(TAG_COUNT) 

help:	## show this help
	@grep -E '^[a-zA-Z0-9_-]+:.*?## .*$$' $(MAKEFILE_LIST) | awk 'BEGIN {FS = ":.*?## "}; {printf "%-10s %s\n", $$1, $$2}'


publish: ## ko build
publish: build_up
	az acr login -n $(ACR_NAME)
	env KO_DOCKER_REPO=$(KO_DOCKER_REPO) ko build --base-import-paths -t $(TAG_VERSION)-$(TAG_COUNT) -t latest ./cmd/server

local: ## ko build (ko.local)
	$(MAKE) publish KO_DOCKER_REPO=ko.local

lint:	## run buf lint
lint:	lint-golangci lint-buf

format:	## run buf/Go format
	buf format -w
	go fmt ./...

generate:	## run buf generate
	rm -rf pkg/gen/buf/*
	buf generate

## internal targets

# build number を加算する
build_up:
	@if [[ ! -f $(BUILD_NUMBER_FILE) ]] ; then echo 0 > $(BUILD_NUMBER_FILE); fi
	@echo $$(($$(cat $(BUILD_NUMBER_FILE)) + 1)) > $(BUILD_NUMBER_FILE)

lint-golangci:	# run golangci-lint
	docker run --rm -v $$(pwd):/app -v ~/.cache/golangci-lint/v1.58.0:/root/.cache -w /app golangci/golangci-lint:v1.58.0 golangci-lint run

lint-buf:	# run buf lint
	buf lint

## infra

deploy:	## deploy infra
	env IMAGE_NAME=$(IMAGE_NAME) az deployment sub create --location $(AZURE_LOCATION) --template-file ./infra/main.bicep --parameters infra/main.bicepparam --name $(TAG_VERSION)-$(TAG_COUNT)

provision:	## provision infra
	azd priovision

infra-refresh:	## azd 設定をazure環境に同期
	azd env refresh 
	azd env get-values | tr -d '"' > .env

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

go.work: ## 開発用のgo.workを作成
	go work init ./cmd/client ./cmd/server ./pkg

