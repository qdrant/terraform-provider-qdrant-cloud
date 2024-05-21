default: test

NAME=qdrant-cloud
BINARY=terraform-provider-${NAME}
VERSION=1.0
HOSTNAME=cloud.qdrant.io
OS_ARCH=$(shell uname -s | tr '[:upper:]' '[:lower:]')_$(shell uname -m)
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
	go build -o ${BINARY}

install: build
	mkdir -p ~/.terraform.d/plugins/${HOSTNAME}/${NAMESPACE}/${NAME}/${VERSION}/${OS_ARCH}
	cp ${BINARY} ~/.terraform.d/plugins/${HOSTNAME}/${NAMESPACE}/${NAME}/${VERSION}/${OS_ARCH}
