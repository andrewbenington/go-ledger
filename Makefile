.PHONY: start
start:
	go run .

.PHONY: build
build:
	go build -o bin/go-ledger .

.PHONY: lint
lint:
	golangci-lint run

.PHONY: clean
clean:
	rm -f ./*.xlsx

.PHONY: deps
deps:
	go get ./...