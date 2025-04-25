SHELL := bash

GO_LD_FLAGS ?= "-w -s"

# Default target to build the project
default: build

lint: ## Run linter
	golangci-lint run -v --timeout 300s

test: ## Run unit-tests
	go test -v -race -timeout 300s -coverprofile=coverprofile.out -coverpkg "./pkg/..." ./...

build: ## Build the project with optional GOOS and GOARCH
	( GOOS="$(GOOS)" GOARCH="$(GOARCH)" GO_BUILD_FLAGS=$(GO_BUILD_FLAGS) hack/build.sh )

build-all: ## Build the project for all supported platforms
	hack/build-all.sh

coverprofile: ## Generate coverage profile
	hack/coverprofile.sh

fmt: ## Format Go code
	gofmt -s -w ./cmd ./pkg ./tests

validate: ## Validate that the generated files are up to date
	hack/validate.sh

assets: ## Generate assets for the project
	hack/generate-assets.sh

update-deps: ## Update project dependencies
	hack/update-deps.sh

generate: ## Run Go generate for the project
	go generate ./...

lint-metainfo: ## Lint the metainfo file for Flatpak
	flatpak run --command=flatpak-builder-lint org.flatpak.Builder appstream io.github.heathcliff26.go-minesweeper.metainfo.xml

gosec: ## Scan code for vulnerabilities using gosec
	gosec ./...

clean: ## Clean up build artifacts and temporary files
	hack/clean.sh

help: ## Show this help message
	@echo "Available targets:"
	@grep -E '^[a-zA-Z_-]+:.*?##' $(MAKEFILE_LIST) | awk 'BEGIN {FS = ":.*?## "}; {printf "%-20s %s\n", $$1, $$2}'
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
