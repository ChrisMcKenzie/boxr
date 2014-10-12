SELFPKG := github.com/Secret-Ironman/boxr
VERSION := 0.0.1a
SHA := $(shell git rev-parse --short HEAD)
BRANCH := $(shell git rev-parse --abbrev-ref HEAD)

deps:
	go get -u -t -v ./...

godep:
	go get github.com/tools/godep

deps-save: 
	godep save -r ./...

build: 
	go build -o bin/boxr -ldflags "-X main.version $(VERSION)dev-$(SHA)" $(SELFPKG)/cli
	go build -o bin/boxrd -ldflags "-X main.version $(VERSION)dev-$(SHA)" $(SELFPKG)/service

test: godep
	godep go test -v ./...

run: 
	bin/boxrd -port=:3000