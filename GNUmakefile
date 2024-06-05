default: test

NAME=qdrant-cloud
BINARY=terraform-provider-${NAME}
VERSION=1.0
HOSTNAME=cloud.qdrant.io
OS=$(shell uname -s | tr '[:upper:]' '[:lower:]')
ARCH=$(shell uname -m)
OS_ARCH=${OS}_${ARCH}
NAMESPACE=$(shell whoami)
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

generate-client:
	cd internal && swagger-codegen generate -i ./spec.json -l go --output client --additional-properties packageName=cloud

build:
	CGO_ENABLED=0 ~/go/bin/gox \
		-osarch="linux/amd64 linux/arm linux/arm64 darwin/amd64 darwin/arm64 windows/amd64" \
		-output="bin/{{.OS}}/{{.Arch}}/${BINARY}_v$(VERSION)" \
		-tags="netgo" \
		./...

.PHONY: update-go-client
update-go-client:
	rm -r ./go-client-programmatic-access
	mkdir ./go-client-programmatic-access
	cp -R -v ../qdrant-cloud-cluster-api/pypi/go-client-programmatic-access/* ./go-client-programmatic-access

install: build
	mkdir -p ~/.terraform.d/plugins/${HOSTNAME}/${NAMESPACE}/${NAME}/${VERSION}/${OS_ARCH}
	cp bin/${OS}/${ARCH}/${BINARY}_v$(VERSION) ~/.terraform.d/plugins/${HOSTNAME}/${NAMESPACE}/${NAME}/${VERSION}/${OS_ARCH}/${BINARY}

.PHONY: generate-help
generate-help:
	go generate ./...

.PHONY: checksum
checksum:
	find bin -type f -exec sha256sum {} \; > checksums.txt
