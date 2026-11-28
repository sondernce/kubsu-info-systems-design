GOLANGCI_LINT ?= golangci-lint

.PHONY: fmt lint test check

fmt:
	$(GOLANGCI_LINT) fmt

lint:
	$(GOLANGCI_LINT) run ./...

test:
	go test ./...

check: lint test
