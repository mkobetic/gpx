default: build

setup:
	@go get
	@go install github.com/benbjohnson/ego/cmd/ego

build: *.go ego.go
	@go build

ego.go: *.ego
	@ego

samples: build
	@rm -f samples/out/*
	@./gpx -o samples/out -vo 0 samples/in/*

.PHONY: setup samples