GOFILES = $(shell find ./cmd/orchard/ -name '*.go' -not -path './vendor/*')
GOPACKAGES = $(shell go list ./...  | grep -v /vendor/)

default: build

workdir:
	mkdir -p workdir

build: workdir/orchard

build-native: $(GOFILES)
	go build -o build/orchard ./cmd/orchard

workdir/orchard: $(GOFILES)
	GOOS=darwin GOARCH=amd64 CGO_ENABLED=0 go build -o build/orchard-mac-amd64 ./cmd/orchard
	GOOS=linux GOARCH=amd64 CGO_ENABLED=0 go build -o build/orchard-linux-amd64 ./cmd/orchard
	GOOS=windows GOARCH=amd64 CGO_ENABLED=0 go build -o build/orchard-windows-amd64 ./cmd/orchard

test: test-all

test-all:
	@go test -race -coverprofile=coverage.txt -covermode=atomic $(GOPACKAGES)
