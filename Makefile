VERSION = $(shell git rev-parse --short HEAD)

all: help

help:
	@echo
	@echo "Commands"
	@echo "========"
	@echo
	@sed -n '/^[a-zA-Z0-9_-]*:/s/:.*//p' < Makefile | grep -v -E 'default|help.*' | sort


gen-proto:
	protoc \
		-I $(HOME)/Downloads/protobuf-3.19.4/src/ \
		-I .  \
		--go_out=. --go_opt=paths=source_relative \
		--go-grpc_out=. --go-grpc_opt=paths=source_relative \
		pkg/goe/goe.proto


build-app-%:
	mkdir -p bin
	CGO_ENABLED=1 go build -o bin/${*} github.com/husio/goe-demo/app/${*}

build-image-%:
	docker build -t "${*}:${VERSION}" -t "${*}:latest" -f app/${*}/Dockerfile .


build-images: build-image-consumer build-image-producer


latest-store-entries:
	@watch 'docker-compose run redis redis-cli -h redis --raw LRANGE randomer 0 10'
