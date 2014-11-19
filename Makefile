SELFPKG := github.com/Secret-Ironman/boxr
VERSION := 0.0.1a
SHA := $(shell git rev-parse --short HEAD)
BRANCH := $(shell git rev-parse --abbrev-ref HEAD)

deps: godep
	godep restore

godep:
	go get github.com/tools/godep

deps-save: 
	godep save -r ./...

build:
	godep go build -o bin/boxr -ldflags "-X main.version $(VERSION)dev-$(SHA)" $(SELFPKG)/main.go

test:
	godep go test -v ./...

run: build
	bin/boxr s 