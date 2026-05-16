.PHONY: build test

build:
	go build -o build/vx cmd/vx/*.go

test:
	go test ./... -count=1 -timeout=60s
