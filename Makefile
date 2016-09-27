default: build

setup:
	@go get
	@go get github.com/benbjohnson/ego/cmd/ego

build: *.go ego.go
	@go build

ego.go: *.ego
	@ego -package=main *.ego

samples:
	@rm samples/out/*
	@./gpx -o samples/out samples/in/*

.PHONY: setup samples