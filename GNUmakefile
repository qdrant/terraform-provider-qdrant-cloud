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
