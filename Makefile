.PHONY: deps-install test build docker-build

BUILD_EXECUTABLE := smsender

all: build

deps-install:
	go get -v -u github.com/Masterminds/glide
	glide install

test:
	@go test -race -v $(shell go list ./... | grep -v vendor)

build:
	go build -o ./bin/$(BUILD_EXECUTABLE)

docker-build:
	GOOS=linux GOARCH=amd64 go build -o ./bin/$(BUILD_EXECUTABLE)