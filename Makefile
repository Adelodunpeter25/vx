.PHONY: build test

build:
	go build -o build/vx cmd/vx/*.go

test: build
	./build/vx
