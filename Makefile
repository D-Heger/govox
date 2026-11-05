VERSION := $(shell git describe --tags --always --dirty 2>/dev/null || echo dev)
BIN := bin/govox
SRC := $(shell find cmd internal -name '*.go') go.mod $(wildcard go.sum)

.PHONY: count build run test package clean

count:
	find . -name '*.go' | xargs wc -l

build: $(BIN)
$(BIN): $(SRC)
	go mod tidy
	CGO_ENABLED=0 go build -trimpath -o $@ ./cmd/govox

run: $(BIN)
	./$(BIN)

test:
	gofmt -l . && go vet ./... && go test -race ./...

package:
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o dist/govox-linux-amd64 ./cmd/govox
	CGO_ENABLED=0 GOOS=darwin GOARCH=arm64 go build -o dist/govox-macos-arm64 ./cmd/govox
	CGO_ENABLED=0 GOOS=windows GOARCH=amd64 go build -o dist/govox-windows-amd64.exe ./cmd/govox