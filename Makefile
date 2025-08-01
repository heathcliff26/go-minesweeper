SHELL := bash

GO_LD_FLAGS ?= "-w -s"

# Default target to build the project
default: build

# Run linter
lint:
	golangci-lint run -v --timeout 300s

# Run unit-tests
test:
	go test -v -race -timeout 300s -coverprofile=coverprofile.out -coverpkg "./pkg/..." ./...

# Build the binary
build: tools
	"$(shell pwd)/bin/fyne" build -o "$(shell pwd)/bin/go-minesweeper" -release

# Build the project for all supported platforms
build-all:
	hack/build-all.sh

# Generate coverage profile
coverprofile:
	hack/coverprofile.sh

# Format Go code
fmt:
	gofmt -s -w ./cmd ./pkg ./tests

# Validate that the generated files are up to date
validate:
	hack/validate.sh

# Validate the appstream metainfo file
validate-metainfo:
	appstreamcli validate io.github.heathcliff26.go-minesweeper.metainfo.xml

# Generate assets for the project
assets:
	hack/generate-assets.sh

# Update project dependencies
update-deps:
	hack/update-deps.sh

# Run Go generate for the project
generate:
	go generate ./...

# Scan code for vulnerabilities using gosec
gosec:
	gosec ./...

# Clean up build artifacts and temporary files
clean:
	hack/clean.sh

# Install the tools required for building the app
tools:
	GOBIN="$(shell pwd)/bin" go install tool

# Show this help message
help:
	@echo "Available targets:"
	@echo ""
	@awk '/^#/{c=substr($$0,3);next}c&&/^[[:alpha:]][[:alnum:]_-]+:/{print substr($$1,1,index($$1,":")),c}1{c=0}' $(MAKEFILE_LIST) | column -s: -t
	@echo ""
	@echo "Run 'make <target>' to execute a specific target."

.PHONY: \
	default \
	build \
	test \
	lint \
	coverprofile \
	fmt \
	validate \
	validate-metainfo \
	assets \
	update-deps \
	generate \
	gosec \
	clean \
	tools \
	help \
	$(NULL)
