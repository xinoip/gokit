include .env
export

.PHONY: *

all: test

tidy:
	go mod tidy

build: gen tidy
	go build -o server.out .

lint:
	golangci-lint run

gen:
	go generate ./...
	go tool sqlc -f sqlc.yaml generate

clean:
	rm -rf server.out internal/gen
	go clean -testcache

test: build
	go test -coverprofile=coverage.out -race ./...

test-debug:
	go test -v -parallel=1 ./...

doc:
	pkgsite -open

run: build
	./server.out serve

migrate: build
	./server.out migrate

openapi: build
	./server.out openapi

up-dependencies:
	docker-compose up -d
