# Copyright (C) 2020 ConsenSys Software Inc.

# Build the Filecoin Retrieval Gateway

# Usage:
#   [VERSION=v3] [REGISTRY="gcr.io/google_containers"] make build
VERSION?=v1
REGISTRY?=

release: clean build push clean

# builds a docker image that builds the app and packages it into a minimal docker image
build:
	docker build -t ${REGISTRY}fc-retrieval-client-builder .
	docker run --rm ${REGISTRY}fc-retrieval-client-builder | docker build --pull -t "${REGISTRY}fc-retrieval-client:${VERSION}" -

# push the image to an registry
push:
#	gcloud docker -- push ${REGISTRY}/fc-retrieval-client:${VERSION}

utest:
	go test ./...

# remove previous images and containers
clean:
#	rm -f /etc/client/
	docker rm -f ${REGISTRY}fc-retrieval-client-builder 2> /dev/null || true
	docker rmi -f ${REGISTRY}fc-retrieval-client-builder || true
	docker rmi -f "${REGISTRY}fc-retrieval-client:${VERSION}" || true

.PHONY: release clean build push

