GOLANGCI_VERSION ?= v1.61.0
GOVULNCHECK_VERSION ?= v1.1.3
GOFUMPT_VERSION ?= v0.7.0

ROOT_DIR:=$(shell dirname $(realpath $(firstword $(MAKEFILE_LIST))))

.PHONY: default
default: test

.PHONY: install-govulncheck
install-govulncheck:
	@GOBIN=$(ROOT_DIR)/bin go install golang.org/x/vuln/cmd/govulncheck@$(GOVULNCHECK_VERSION)

.PHONY: install-golangci
install-golangci:
	@test -f ./bin/golangci-lint || curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | sh -s -- $(GOLANGCI_VERSION)

.PHONY: install-gofumpt
install-gofumpt:
	@GOBIN=$(ROOT_DIR)/bin go install mvdan.cc/gofumpt@$(GOFUMPT_VERSION)

.PHONY: go-lint
go-lint: install-govulncheck install-golangci
	@./bin/golangci-lint run

.PHONY: go-fumpt
go-fumpt: install-gofumpt
	@./bin/gofumpt -l -w .

.PHONY: govulncheck
govulncheck: install-govulncheck
	@./bin/govulncheck ./...

.PHONY: lint
lint: go-lint go-fumpt govulncheck

.PHONY: test
test:
	go test -cover -race -count=100 -short -race ./...

.PHONY: bench
bench:
	go test -bench . -test.benchmem -test.count 2 -test.benchtime 2s
