.PHONY: build clean install test

build:
	go build -o vx cmd/vx/*.go

clean:
	rm -f vx

install: build
	cp vx /usr/local/bin/

test:
	go test ./...

run: build
	./vx

.DEFAULT_GOAL := build
