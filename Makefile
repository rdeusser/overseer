NAME=overseer
ARCH=$(shell uname -m)
VERSION=0.0.1

all: fmt deps test prod

fmt:
	go fmt `go list ./...`

deps:
	go get -t -v ./...

test: deps
	go tool vet .
	golint ./...
	go test -v -race ./...

prod: fmt deps test
	@M_PROD=1 BIN_VERSION=$(VERSION) sh -c "'$(CURDIR)/scripts/build.sh'"

dev: fmt deps test
	@M_DEV=1 BIN_VERSION=$(VERSION) sh -c "'$(CURDIR)/scripts/build.sh'"

release: prod
	rm -rf release && mkdir release
	tar -cvzf release/$(NAME)_$(VERSION)_linux_$(ARCH).tar.gz -C build/linux $(NAME)
	tar -cvzf release/$(NAME)_$(VERSION)_darwin_$(ARCH).tar.gz -C build/darwin $(NAME)

.PHONY: all fmt deps test prod dev release
