.PHONY: build test lint run

build:
	go build -o bin/bernstein ./cmd/bernstein

test:
	go test -v ./...

lint:
	golangci-lint run

run: build
	./bin/bernstein

format:
	goimports -w .
	gofmt -s -w .

