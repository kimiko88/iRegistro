
# Makefile for iRegistro

.PHONY: all build test lint run docker-build docker-run clean

PROJECT_NAME := iRegistro
BINARY_NAME := iRegistro
DOCKER_IMAGE := iRegistro

all: build

build:
	go build -o bin/$(BINARY_NAME) cmd/api/main.go

test:
	go test -v ./... -cover

lint:
	golangci-lint run

run:
	go run cmd/api/main.go

docker-build:
	docker build -t $(DOCKER_IMAGE) .

docker-run:
	docker run -p 8080:8080 $(DOCKER_IMAGE)

clean:
	rm -rf bin/
