GO ?= go

build:
	$(GO) build

test:
	$(GO) test

.PHONY: build

extras:
	$(GO) get golang.org/x/tools/cmd/vet
	$(GO) get github.com/golang/lint/golint
	$(GO) vet
	golint

.PHONY: extras

check: extras test

