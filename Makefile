.PHONY: install build test lint bench

default: build test

install:
	go mod download

build:
	go build ./...

test:
	go test ./...

lint:
	golangci-lint run ./...

bench:
	go test -bench=. -benchmem ./...
