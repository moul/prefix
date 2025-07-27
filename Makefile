.PHONY: all build test clean install lint

all: test build

build:
	go build -o bin/prefix ./cmd/prefix

test:
	go test -v -race -coverprofile=coverage.out ./...

coverage: test
	go tool cover -html=coverage.out -o coverage.html

clean:
	rm -rf bin/ coverage.out coverage.html

install:
	go install ./cmd/prefix

lint:
	golangci-lint run

docker:
	docker build -t prefix .

.DEFAULT_GOAL := all