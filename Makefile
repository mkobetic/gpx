BUILT := $(shell date -u '+%Y-%m-%d %I:%M:%S')
TAG := $(shell git tag --points-at HEAD)
COMMIT := $(shell git rev-parse --short HEAD)
BRANCH := $(shell git rev-parse --abbrev-ref HEAD)
GO_VERSION := $(shell go version)
LDFLAGS += -X "main.Built=$(BUILT)"
LDFLAGS += -X "main.Commit=$(COMMIT)/$(TAG)"
LDFLAGS += -X "main.Branch=$(BRANCH)"
LDFLAGS += -X "main.GoVersion=$(GO_VERSION)"
BUILD := go build -ldflags '$(LDFLAGS)'

default: build

setup:
	@go get
	@go install github.com/benbjohnson/ego/cmd/ego

build: *.go *.ego
	ego
	$(BUILD) .

install: *.go *.ego
	ego
	go install -ldflags '$(LDFLAGS)' .

ego: *.ego
	@ego

samples: build
	@rm -f samples/out/*
	@./gpx -o samples/out -wd NW -vo 0 samples/in/*

dist:
	GOOS=darwin GOARCH=arm64 $(BUILD) -o dist/ .
	tar zcf dist/gpx-osx-arm64.tgz -C dist/ gpx
	GOOS=darwin GOARCH=amd64 $(BUILD) -o dist/ .
	tar zcf dist/gpx-osx-amd64.tgz -C dist/ gpx
	GOOS=linux GOARCH=arm64 $(BUILD) -o dist/ .
	tar zcf dist/gpx-linux-arm64.tgz -C dist/ gpx
	GOOS=linux GOARCH=amd64 $(BUILD) -o dist/ .
	tar zcf dist/gpx-linux-amd64.tgz -C dist/ gpx
	GOOS=windows GOARCH=amd64 $(BUILD) -o dist/ .
	tar zcf dist/gpx-windows-amd64.tgz -C dist/ gpx.exe

test:
	@go test .

.PHONY: setup samples dist test