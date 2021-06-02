# Copyright (C) 2020 ConsenSys Software Inc.

# Build the Filecoin Retrieval Gateway

# Usage:
#   [VERSION=v3] [REGISTRY="gcr.io/google_containers"] make build
VERSION?=dev
IMAGE?=consensys/fc-retrieval-itest
COMPOSE_FILE?=docker-compose.yml

# Always assume these targets are out of date.
.PHONY: clean itest itest-dev utest build release push detectmisconfig


# This target (the first target in the build file) is the one that is executed if no 
# command line args are specified.
default: clean utest build tag

# builds a docker image that builds the app and packages it into a minimal docker image
build:
	docker build -t ${IMAGE}:${VERSION} .

lotusbase:
	docker build -t consensys/lotus-base lotus/lotus-base
lotusdaemon:
	docker build -t consensys/lotus-daemon lotus/lotus-daemon
lotusfullnode:
	docker build -t consensys/lotus-full-node lotus/lotus-full-node

buildlocal:
	cd ..; docker build -f ./fc-retrieval-itest/Dockerfile.local -t ${IMAGE}:${VERSION} .

# push the image to an registry
push:
	cd scripts; bash push.sh ${VERSION} ${IMAGE}:${VERSION}

tag:
	cd scripts; bash tag.sh ${VERSION} ${IMAGE}:${VERSION}

uselocal:
	echo "replace github.com/ConsenSys/fc-retrieval-common => ../fc-retrieval-common" >> go.mod
	echo "replace github.com/ConsenSys/fc-retrieval-client => ../fc-retrieval-client" >> go.mod
	echo "replace github.com/ConsenSys/fc-retrieval-gateway-admin => ../fc-retrieval-gateway-admin" >> go.mod
	echo "replace github.com/ConsenSys/fc-retrieval-provider-admin => ../fc-retrieval-provider-admin" >> go.mod
	go mod tidy

useremote:
	cd scripts; bash use-remote-repos.sh

detectmisconfig:
	cd scripts; bash detect-pkg-misconfig.sh

# Local build: make sure the test code compiles. 
lbuild:
	go test -c github.com/ConsenSys/fc-retrieval-itest/pkg/client-gateway
	go test -c github.com/ConsenSys/fc-retrieval-itest/pkg/client-init
	go test -c github.com/ConsenSys/fc-retrieval-itest/pkg/poc1
	go test -c github.com/ConsenSys/fc-retrieval-itest/pkg/poc2
	go test -c github.com/ConsenSys/fc-retrieval-itest/pkg/provider-admin
	go test -c github.com/ConsenSys/fc-retrieval-itest/pkg/lotus
	

# Test all packages, skipping lotus-full-node on CircleCI
itestlocal: 
	go test -v -p=1 --count=1 ./pkg/lotus-full-node

setup-env-localtesting:
	cd scripts; bash setup-env.sh

# Version that can be run on a desktop computer or in Circle CI.
# Itest run from a container.
# Run the gateway(s), provider(s), and register services in Docker. Run the 
# tests locally. Dump the go.mod file so that the precise versions of 
# Client and Gateway Admin library are recorded. 
itestdocker:
	bash ./scripts/itestdocker.sh
	
# This is the previous methodology, where the integration tests were in 
# a Docker container.
#
#	docker-compose down
#	docker-compose up --abort-on-container-exit --exit-code-from itest

# Dump network config:
#		echo NETWORK CONFIG *********************************************; \
#		docker network inspect shared; \


clean:
	docker rm -f ${IMAGE}:${VERSION} 2> /dev/null || true
	docker rmi -f ${IMAGE}:${VERSION} || true

check-modules:
	./scripts/check-modules/check-modules

check-main-modules:
	./scripts/check-main-modules/check-main-modules

dev:
	./scripts/make-env-file/make-env-file -source=.env.example -dest=.env
	docker-compose -f docker-compose.dev.yml up $(arg)

test:
	for number in 1 2 3 4 ; do \
		echo $$number ; \
	done