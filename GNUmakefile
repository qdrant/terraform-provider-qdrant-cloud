SHELL := /bin/sh

default: test

NAME=qdrant-cloud
BINARY=terraform-provider-${NAME}
VERSION=1.0
OS=$(shell uname -s | tr '[:upper:]' '[:lower:]')
ARCH=$(shell uname -m)
OS_ARCH=${OS}_${ARCH}
NAMESPACE=local
QDRANT_CLOUD_API_KEY ?=
QDRANT_CLOUD_ACCOUNT_ID ?=

# Get the latest tag
LATEST_TAG := $(shell git describe --tags `git rev-list --tags --max-count=1`)

# Use the latest tag in a target
print-latest-tag:
	@echo "The latest tag is $(LATEST_TAG)"

.PHONY: test
# Run acceptance tests
test: test.validate_env_vars test.acceptance

.PHONY: test.validate_env_vars
# Validate environment variables required for acceptance tests
test.validate_env_vars:
	@echo "🔍 Validating test environment..."; \
	missing=""; \
	[ -z "$$QDRANT_CLOUD_API_KEY" ] && missing="$$missing QDRANT_CLOUD_API_KEY"; \
	[ -z "$$QDRANT_CLOUD_ACCOUNT_ID" ] && missing="$$missing QDRANT_CLOUD_ACCOUNT_ID"; \
	[ -z "$$QDRANT_CLOUD_API_URL" ] && missing="$$missing QDRANT_CLOUD_API_URL"; \
	if [ -n "$$missing" ]; then \
		echo "❌ Missing required environment variables:$$missing"; \
		exit 1; \
	else \
		echo "✅ All required environment variables are set."; \
	fi


.PHONY: test.unit
test.unit:
	@set -e; \
	echo "🧪 Running unit tests..."; \
	tests_unit="$$(grep -Rho '^func Test[A-Za-z0-9_]*' ./qdrant | awk '{print $$2}' | grep -v '^TestAcc')"; \
	skips="$(SKIPPED_UNIT_TESTS)"; \
	if [ -z "$$tests_unit" ]; then \
		echo "No non-acceptance tests found."; \
		exit 0; \
	fi; \
	for t in $$tests_unit; do \
		if [ -n "$$skips" ] && echo "$$skips" | tr ' ' '\n' | grep -qx "$$t"; then \
			echo "\n⏭  Skipping $$t"; \
			continue; \
		fi; \
		echo ""; echo "==== ▶ $$t ===="; \
		if [ "$(DRYRUN)" = "1" ]; then \
			echo "go test -count=1 -v ./qdrant -run '^$$t$$'"; \
		else \
			go test -count=1 -v ./qdrant -run "^$$t$$" || { echo "✖ $$t FAILED"; exit 1; }; \
			echo "✓ $$t passed"; \
		fi; \
	done; \
	echo ""; echo "✅ Unit tests completed."


.PHONY: test.acceptance
test.acceptance:
	@set -e; \
	echo "🚀 Running acceptance tests..."; \
	tests="$$(grep -Rho '^func TestAcc[A-Za-z0-9_]*' . | awk '{print $$2}')"; \
	skips="$(SKIPPED_TESTS)"; \
	for testname in $$tests; do \
		if [ -n "$$skips" ] && echo "$$skips" | tr ' ' '\n' | grep -qx "$$testname"; then \
			echo "\n⏭  Skipping $$testname"; \
			continue; \
		fi; \
		echo ""; \
		echo "==== ▶ $$testname ===="; \
		if [ "$(DRYRUN)" = "1" ]; then \
			echo "TF_ACC=1 QDRANT_CLOUD_API_KEY=\"$(QDRANT_CLOUD_API_KEY)\" QDRANT_CLOUD_ACCOUNT_ID=\"$(QDRANT_CLOUD_ACCOUNT_ID)\" QDRANT_CLOUD_API_URL=\"$(QDRANT_CLOUD_API_URL)\" go test -count=1 -v ./... -run '^$$testname$$'"; \
		else \
			TF_ACC=1 \
			QDRANT_CLOUD_API_KEY="$(QDRANT_CLOUD_API_KEY)" \
			QDRANT_CLOUD_ACCOUNT_ID="$(QDRANT_CLOUD_ACCOUNT_ID)" \
			QDRANT_CLOUD_API_URL="$(QDRANT_CLOUD_API_URL)" \
			go test -count=1 -v ./... -run "^$$testname$$" || { echo "✖ $$testname FAILED"; exit 1; }; \
			echo "✓ $$testname passed"; \
		fi; \
	done; \
	echo ""; \
	echo "✅ All acceptance tests completed."

requirements:
	go install github.com/goreleaser/goreleaser/v2@latest
	go install github.com/mitchellh/gox@latest

build: requirements
	goreleaser release --snapshot --clean

install: build
	mkdir -p ~/.terraform.d/plugins/${NAMESPACE}/${NAME}/${NAME}/${VERSION}/${OS_ARCH}
	cp bin/${OS}/${ARCH}/${BINARY}_v$(VERSION) ~/.terraform.d/plugins/${NAMESPACE}/${NAME}/${NAME}/${VERSION}/${OS_ARCH}/${BINARY}

install-in-gobin: build
	mkdir -p ~/.terraform.d/plugins/${NAMESPACE}/${NAME}/${NAME}/${VERSION}/${OS_ARCH}
	cp dist/${BINARY}_${OS}_${ARCH}/${BINARY}_$(LATEST_TAG)-next ${HOME}/go/bin/terraform-provider-qdrant-cloud

.PHONY: generate-help
generate-help:
	go generate ./...

.PHONY: checksum
checksum:
	find bin -type f -exec sha256sum {} \; > checksums.txt

.PHONY: local-build
local-build:
	mkdir -p ~/.terraform.d/plugins/${NAMESPACE}/${NAME}/${NAME}/${VERSION}/${OS_ARCH}
	go build
	mv ${BINARY} ~/.terraform.d/plugins/${NAMESPACE}/${NAME}/${NAME}/${VERSION}/${OS_ARCH}/${BINARY}

.PHONY: lint
lint: bootstrap ## Run project linters
	$(GOLANGCI_LINT) run

##@ Dependencies

## Location to install dependencies to
LOCALBIN ?= $(shell pwd)/bin
$(LOCALBIN):
	mkdir -p $(LOCALBIN)

## Tool Binaries
GOLANGCI_LINT = $(LOCALBIN)/golangci-lint

## Tool Versions
GOLANGCI_LINT_VERSION ?= v2.6.2

.PHONY: bootstrap
bootstrap: install/golangci-lint ## Install required dependencies to work with this project

.PHONY: golangci-lint
install/golangci-lint: $(GOLANGCI_LINT) ## Download golangci-lint locally if necessary.
$(GOLANGCI_LINT): $(LOCALBIN)
	$(call go-install-tool,$(GOLANGCI_LINT),github.com/golangci/golangci-lint/v2/cmd/golangci-lint,$(GOLANGCI_LINT_VERSION))

# copied from kube-builder
# go-install-tool will 'go install' any package with custom target and name of binary, if it doesn't exist
# $1 - target path with name of binary
# $2 - package url which can be installed
# $3 - specific version of package
define go-install-tool
@[ -f "$(1)-$(3)" ] || { \
set -e; \
package=$(2)@$(3) ;\
echo "Downloading $${package}" ;\
rm -f $(1) || true ;\
GOBIN=$(LOCALBIN) go install $${package} ;\
mv $(1) $(1)-$(3) ;\
} ;\
ln -sf $(1)-$(3) $(1)
endef
