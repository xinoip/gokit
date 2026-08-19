.PHONY: *

all: tidy lint lint-example test test-example

tidy:
	go mod tidy
	cd example && go mod tidy

lint:
	golangci-lint run

lint-example:
	cd example && make lint

test:
	go test -coverprofile=coverage.out -race ./...

test-example:
	cd example && make test

clean:
	go clean -testcache

doc:
	pkgsite -open
