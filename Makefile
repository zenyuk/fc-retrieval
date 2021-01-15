# Copyright (C) 2020 ConsenSys Software Inc.

# Build the Filecoin Retrieval Gateway

# Usage:
#   [VERSION=v3] [REGISTRY="gcr.io/google_containers"] make build
VERSION?=v1
REGISTRY?=

# This target (the first target in the build file) is the one that is executed if no 
# command line args are specified.
release: clean utest build

# builds a docker image that builds the app and packages it into a minimal docker image
build:
	docker build -t ${REGISTRY}fc-retrieval-itest:${VERSION} .


# push the image to an registry
push:
#	gcloud docker -- push ${REGISTRY}/fc-retrieval-client:${VERSION}

detectmisconfig:
	cd scripts; bash detect-pkg-misconfig.sh

utest:
	go test ./...

itest:
	docker-compose down
	docker-compose up --abort-on-container-exit --exit-code-from itest


# remove previous images and containers
clean:
	docker rm -f ${REGISTRY}fc-retrieval-itest-builder 2> /dev/null || true
	docker rmi -f ${REGISTRY}fc-retrieval-itest-builder || true
	docker rmi -f "${REGISTRY}fc-retrieval-itest:${VERSION}" || true

# Alays assume these targets are out of date.
.PHONY: clean itest utest build release push detectmisconfig

