.DEFAULT_GOAL := run

fmt: 
	go fmt ./...

lint: fmt
	golint ./...

vet: fmt
	go vet ./...

build: vet
	go build -o bin/main cmd/api/main.go

run: build
	go run cmd/api/main.go

clean:
	go clean
	rm bin/main

.PHONY: fmt lint vet build run clean