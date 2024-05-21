HOSTNAME=qdrant.com
NAMESPACE=com
NAME=qdrant-cloud
BINARY=terraform-provider-${NAME}
VERSION=0.2
OS_ARCH=darwin_amd64

default: install

clean:
	rm -Rf data_profile/
	rm -Rf qdrantcloud/*
	rm -Rf client/*
	rm -Rf models/*
	rm -f  main.go

build: clean
	/Users/p0123n/go/bin/swagger generate client -f ./spec-v2.json --template-dir templates -C config.yml > swagrun.log

compile:
	go mod tidy
	go build -o ${BINARY}

install: build
	#mkdir -p ~/.terraform.d/plugins/${HOSTNAME}/${NAMESPACE}/${NAME}/${VERSION}/${OS_ARCH}
	#cp ${BINARY} ~/.terraform.d/plugins/${HOSTNAME}/${NAMESPACE}/${NAME}/${VERSION}/${OS_ARCH}

nogen: 
	go build -o ${BINARY}
	#mkdir -p ~/.terraform.d/plugins/${HOSTNAME}/${NAMESPACE}/${NAME}/${VERSION}/${OS_ARCH}
	#mv ${BINARY} ~/.terraform.d/plugins/${HOSTNAME}/${NAMESPACE}/${NAME}/${VERSION}/${OS_ARCH}