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
test:
ifndef QDRANT_CLOUD_API_KEY
	$(error QDRANT_CLOUD_API_KEY is not set)
endif
ifndef QDRANT_CLOUD_ACCOUNT_ID
	$(error QDRANT_CLOUD_ACCOUNT_ID is not set)
endif
	TF_ACC=1 go test ./qdrant/... -v $(TESTARGS) -timeout 120m

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
	cp ${BINARY} ~/.terraform.d/plugins/${NAMESPACE}/${NAME}/${NAME}/${VERSION}/${OS_ARCH}/${BINARY}

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
GOLANGCI_LINT_VERSION ?= v2.0.1

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