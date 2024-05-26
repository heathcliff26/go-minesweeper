SHELL := bash

GO_LD_FLAGS ?= "-w -s"

default: build

lint:
	golangci-lint run -v --timeout 300s

test:
	go test -v -race -timeout 10m ./...

build:
	( GOOS="$(GOOS)" GOARCH="$(GOARCH)" GO_BUILD_FLAGS=$(GO_BUILD_FLAGS) hack/build.sh )

build-all:
	hack/build-all.sh

coverprofile:
	hack/coverprofile.sh

assets:
	hack/generate-assets.sh

dependencies:
	hack/update-deps.sh

generate:
	go generate ./...

.PHONY: \
	default \
	build \
	test \
	lint \
	coverprofile \
	assets \
	dependencies \
	generate \
	$(NULL)
