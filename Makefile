.PHONY: default clean install lint test assets build binaries test-release release

TAG_NAME := $(shell git describe --abbrev=0 --tags)
SHA := $(shell git rev-parse --short HEAD)
VERSION := $(if $(TAG_NAME),$(TAG_NAME),$(SHA))
BUILD_DATE := $(shell date -u '+%Y-%m-%d_%H:%M:%S')

IMAGE := andig/evcc
ALPINE := 3.12
TARGETS := arm.v6,arm.v8,amd64

default: clean install assets lint test build

clean:
	rm -rf dist/

install:
	go install github.com/mjibson/esc
	go install github.com/golang/mock/mockgen

lint:
	golangci-lint run

test:
	@echo "Running testsuite"
	go test ./...

assets:
	@echo "Generating embedded assets"
	go generate ./...

build:
	@echo Version: $(VERSION) $(BUILD_DATE)
	go build -v -tags=release -ldflags '-X "github.com/andig/evcc/server.Version=${VERSION}" -X "github.com/andig/evcc/server.Commit=${SHA}"'

release-test:
	goreleaser --snapshot --skip-publish --rm-dist

release:
	goreleaser --rm-dist
