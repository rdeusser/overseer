.PHONY: all fmt deps test bin dev

all: fmt deps test

fmt:
	go fmt `go list ./...`

deps:
	go get -t -v ./...

test: deps
	go tool vet .
	go test -v -race ./...

bin: fmt deps test
	@M_RELEASE=1 sh -c "'$(CURDIR)/scripts/build.sh'"

dev: fmt deps test
	@M_DEV=1 sh -c "'$(CURDIR)/scripts/build.sh'"
