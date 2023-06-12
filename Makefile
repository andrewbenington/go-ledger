.PHONY: start
start:
	go run .

.PHONY: build
build:
	go build -o bin/go-ledger .

.PHONY: lint
lint:
	golangci-lint run

.PHONY: deps
deps:
	go get ./...