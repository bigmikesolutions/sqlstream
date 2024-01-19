GOLANGCI_VERSION ?= v1.54.2
GOVULNCHECK_VERSION ?= v1.0.1

ROOT_DIR:=$(shell dirname $(realpath $(firstword $(MAKEFILE_LIST))))

.PHONY: default
default: test

.PHONY: install-govulncheck
install-govulncheck:
	@GOBIN=$(ROOT_DIR)/bin go install golang.org/x/vuln/cmd/govulncheck@$(GOVULNCHECK_VERSION)

.PHONY: install-golangci
install-golangci:
	@test -f ./bin/golangci-lint || curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | sh -s -- $(GOLANGCI_VERSION)

.PHONY: go-lint
go-lint: install-govulncheck install-golangci
	@./bin/govulncheck ./...
	@./bin/golangci-lint run

.PHONY: lint
lint: go-lint

.PHONY: test
test:
	go test -cover -race -count=100 -short -race ./...

