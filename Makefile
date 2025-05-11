SHELL := bash

GO_LD_FLAGS ?= "-w -s"

# Default target to build the project
default: build

# Run linter
lint:
	golangci-lint run -v --timeout 300s

# Run unit-tests
test:
	go test -v -race -timeout 420s -coverprofile=coverprofile.out -coverpkg "./pkg/..." ./...

# Build the project with optional GOOS and GOARCH
build:
	( GOOS="$(GOOS)" GOARCH="$(GOARCH)" GO_BUILD_FLAGS=$(GO_BUILD_FLAGS) hack/build.sh )

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

# Generate assets for the project
assets:
	hack/generate-assets.sh

# Update project dependencies
update-deps:
	hack/update-deps.sh

# Run Go generate for the project
generate:
	go generate ./...

# Lint the metainfo file for Flatpak
lint-metainfo:
	flatpak run --command=flatpak-builder-lint org.flatpak.Builder appstream io.github.heathcliff26.go-minesweeper.metainfo.xml

# Scan code for vulnerabilities using gosec
gosec:
	gosec ./...

# Clean up build artifacts and temporary files
clean:
	hack/clean.sh

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
	assets \
	update-deps \
	generate \
	lint-metainfo \
	gosec \
	clean \
	help \
	$(NULL)
