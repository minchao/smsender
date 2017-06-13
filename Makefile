.PHONY: deps-install test build build-with-docker docker-build

BUILD_EXECUTABLE := smsender

all: build

deps-install:
	@echo Getting dependencies using Glide
	go get -v -u github.com/Masterminds/glide
	glide install

test:
	@go test -race -v $(shell go list ./... | grep -v vendor)

build: clean
	@echo Building app
	go build -o ./bin/$(BUILD_EXECUTABLE)

clean:
	@echo Cleaning up previous build data
	rm -f ./bin/$(BUILD_EXECUTABLE)

build-with-docker:
	@echo Building app with Docker
	docker run --rm -v $(PWD):/go/src/github.com/minchao/smsender -w /go/src/github.com/minchao/smsender golang sh -c "make deps-install build"

	cd webroot && make build-with-docker

docker-build: build-with-docker
	@echo Building Docker image
	docker build -t minchao/smsender-preview:latest .