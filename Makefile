GOLANGCI_VERSION ?= v2.1.6
GOVULNCHECK_VERSION ?= v1.1.4
GOFUMPT_VERSION ?= v0.8.0

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
	@goimports -l -w .
	@./bin/gofumpt -l -w .

.PHONY: govulncheck
govulncheck: install-govulncheck
	@./bin/govulncheck ./...

.PHONY: lint
lint: go-lint go-fumpt govulncheck

.PHONY: test
test:
	go clean -testcache
	go test -cover -race -count=100 ./...

.PHONY: coverage
coverage:
	go test -v -short -covermode=count -coverprofile=coverage.out ./sql/...
	go tool cover -html=coverage.out -o coverage.html
	go tool cover -func=coverage.out

.PHONY: bench
bench:
	go test -bench . -test.benchmem -test.count 2 -test.benchtime 2s
