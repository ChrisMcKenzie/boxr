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
	go build -o bin/boxr -ldflags "-X main.version $(VERSION)dev-$(SHA)" $(SELFPKG)/cmd/boxr
	go build -o bin/forklift -ldflags "-X main.version $(VERSION)dev-$(SHA)" $(SELFPKG)/cmd/forklift
	go build -o bin/shelfd -ldflags "-X main.version $(VERSION)dev-$(SHA)" $(SELFPKG)/cmd/shelf

test: godep
	godep go test -v ./...

run: 
	bin/boxrd -port=:3000