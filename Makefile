all: fmt lint test

# Initializes the codebase by installing required tools, packages, and git hooks
.PHONY: init
init:
	go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest
	go install mvdan.cc/gofumpt@latest
	go install github.com/daixiang0/gci@latest
	go mod tidy
	echo "#\!/bin/sh\nmake" > .git/hooks/pre-push
	chmod ug+x .git/hooks/pre-push

# Properly formats the project's source files
.PHONY: fmt
fmt:
	gofumpt -l -w .
	gci write --skip-generated -s standard -s default -s "prefix(github.com/goplexhq/needle/internal)" -s blank -s dot .

# Checks the project's source files for linting errors
.PHONY: lint
lint:
	golangci-lint run ./...

# Run all tests in the project
.PHONY: test
test:
	go test ./...
