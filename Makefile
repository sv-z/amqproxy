.PHONY: build
build:
	go build -v ./cmd/proxyserver

.DEFAULT_GOAL := build

.PHONY: test
test:
	go test -v -race -timeout 30s ./...

help:
	./proxyserver -help

run_env:
	docker-compose up -d

stop_env:
	docker-compose down && docker-compose down -v && docker-compose rm -f