.PHONY: build test

build:
	go build -o build/vx cmd/vx/*.go

test:
	go test ./tests/ -count=1 -timeout=60s
