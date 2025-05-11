SHELL := bash

default: help

# Run linter
lint:
	golangci-lint run -v --timeout 300s

# Run unit-tests
test:
	go test -v -race -timeout 300s -coverprofile=coverprofile.out $$(go list ./... | grep -v github.com/heathcliff26/godialog/tests | grep -v github.com/heathcliff26/godialog/examples)

# Run the test app for manual testing
run:
	go build -ldflags="-s" -o "bin/testapp" ./tests/app/...
	bin/testapp

# Generate coverage profile
coverprofile:
	hack/coverprofile.sh

# Format Go code
fmt:
	gofmt -s -w .

# Validate that the generated files are up to date
validate:
	hack/validate.sh

# Update project dependencies
update-deps:
	go mod tidy

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
	lint \
	test \
	run \
	coverprofile \
	fmt \
	validate \
	update-deps \
	gosec \
	clean \
	help \
	$(NULL)
