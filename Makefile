.PHONY: build

build:
	go mod vendor
	go build -ldflags "-s -w " -o ./bin/siproxy ./cmd/main.go