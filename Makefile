.PHONY: *

all: tidy lint test

tidy:
	go mod tidy

lint:
	golangci-lint run

test:
	go test -coverprofile=coverage.out -race ./...

clean:
	go clean -testcache

doc:
	pkgsite -open
