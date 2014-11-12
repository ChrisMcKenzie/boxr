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

build-all: build-boxr build-forklift build-shelf

build-boxr:
	godep go build -o bin/boxr -ldflags "-X main.version $(VERSION)dev-$(SHA)" $(SELFPKG)/cmd/boxr

build-forklift:
	godep go build -o bin/forklift -ldflags "-X main.version $(VERSION)dev-$(SHA)" $(SELFPKG)/cmd/forklift

build-shelf:
	godep go build -o bin/shelf -ldflags "-X main.version $(VERSION)dev-$(SHA)" $(SELFPKG)/cmd/shelf

test: godep
	godep go test -v ./...

run: 
	bin/boxr s &
	bin/forklift &