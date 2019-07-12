.PHONY: check-style test build build-with-docker docker-build

BUILD_EXECUTABLE := smsender
PACKAGES := $(shell go list ./... | grep -v /vendor/)

export GO111MODULE=on

all: build

lint:
	@echo Running lint
	@golangci-lint run -E gofmt ./smsender/...

test:
	@echo Testing
	@go test -race -v $(PACKAGES)

build:
	@echo Building app
	go build -o ./bin/$(BUILD_EXECUTABLE) ./cmd/smsender/main.go

clean:
	@echo Cleaning up previous build data
	rm -f ./bin/$(BUILD_EXECUTABLE)
	rm -rf ./vendor

build-with-docker: clean
	@echo Building app with Docker
	docker run --rm -v $(PWD):/go/src/github.com/minchao/smsender -w /go/src/github.com/minchao/smsender -e GO111MODULE=on golang sh -c "make build"

	cd webroot && make build-with-docker

docker-build: build-with-docker
	@echo Building Docker image
	docker build -t minchao/smsender-preview:latest .
