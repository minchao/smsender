.PHONY: check-style test build build-with-docker docker-build

BUILD_EXECUTABLE := smsender
PACKAGES := $(shell go list ./... | grep -v /vendor/)

export GO111MODULE=on

all: build

vet:
	@echo Running go vet
	@go vet $(PACKAGES)

check-style: vet
	@echo Running go fmt
	$(eval GO_FMT_OUTPUT := $(shell go fmt $(PACKAGES)))
	@echo "$(GO_FMT_OUTPUT)"
	@if [ ! "$(GO_FMT_OUTPUT)" ]; then \
		echo "go fmt success"; \
	else \
		echo "go fmt failure"; \
		exit 1; \
	fi

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
