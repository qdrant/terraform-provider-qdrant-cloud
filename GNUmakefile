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

.PHONY: testacc generate-client
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

generate-client:
	cd internal && swagger-codegen generate -i ./spec.json -l go --output client --additional-properties packageName=cloud

build: requirements
	goreleaser release --snapshot --clean

.PHONY: update-go-client
update-go-client:
	rm -r ./internal/client
	mkdir ./internal/client
	cp -R -v ../qdrant-cloud-cluster-api/pypi/go-client-programmatic-access/* ./internal/client

install: build
	mkdir -p ~/.terraform.d/plugins/${NAMESPACE}/${NAME}/${NAME}/${VERSION}/${OS_ARCH}
	cp bin/${OS}/${ARCH}/${BINARY}_v$(VERSION) ~/.terraform.d/plugins/${NAMESPACE}/${NAME}/${NAME}/${VERSION}/${OS_ARCH}/${BINARY}

.PHONY: generate-help
generate-help:
	go generate ./...

.PHONY: checksum
checksum:
	find bin -type f -exec sha256sum {} \; > checksums.txt
