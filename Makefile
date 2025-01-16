SHELL := bash

GO_LD_FLAGS ?= "-w -s"

default: build

lint:
	golangci-lint run -v --timeout 300s

test:
	go test -v -race -timeout 300s -coverprofile=coverprofile.out ./...

build:
	( GOOS="$(GOOS)" GOARCH="$(GOARCH)" GO_BUILD_FLAGS=$(GO_BUILD_FLAGS) hack/build.sh )

build-all:
	hack/build-all.sh

coverprofile:
	hack/coverprofile.sh

fmt:
	gofmt -s -w ./cmd ./pkg ./tests

validate:
	hack/validate.sh

assets:
	hack/generate-assets.sh

update-deps:
	hack/update-deps.sh

generate:
	go generate ./...

lint-metainfo:
	flatpak run --command=flatpak-builder-lint org.flatpak.Builder appstream io.github.heathcliff26.go-minesweeper.metainfo.xml

clean:
	hack/clean.sh

.PHONY: \
	default \
	build \
	test \
	lint \
	coverprofile \
	fmt \
	validate \
	assets \
	update-deps \
	generate \
	lint-metainfo \
	clean \
	$(NULL)
