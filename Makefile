.PHONY: start
start: clean
	go run .

.PHONY: build
build:
	go build -o bin/go-ledger .

.PHONY: lint
lint:
	golangci-lint run

.PHONY: clean
clean:
	@rm -f ./*.log
	@touch error.log
	@touch go-ledger.log

.PHONY: deps
deps:
	go get ./...