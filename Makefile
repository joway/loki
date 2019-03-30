GOFILES := $(shell find . -name "*.go" -type f -not -path "./vendor/*")
GOFMT ?= gofmt "-s"
VERSION := $(shell cat VERSION)

all: install fmt fmt-check

install:
	@exec ./bin/dep.sh

release:
	@git tag $VERSION
	@git push origin $VERSION

.PHONY: test
test:
	@echo ">> run test"
	@exec env GO111MODULE=on LOKI_ENV='app:xxx*' go test -race -coverprofile=coverage.txt -covermode=atomic -v ./...

fmt:
	@echo ">> formatting code"
	@$(GOFMT) -w $(GOFILES)
