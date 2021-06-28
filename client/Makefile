# Copyright (C) 2020 ConsenSys Software Inc.

# Build the Filecoin Retrieval Gateway

# Usage:
#   [VERSION=v3] [REGISTRY="gcr.io/google_containers"] make build
VERSION?=dev
REGISTRY?=
COV?=80

# This target (the first target in the build file) is the one that is executed if no 
# command line args are specified.
release1: clean useremote utest 

# Use this target if you are using local packages, or if the build is via circle ci, 
# and the go.mod and go.sum file should not be updated
release: clean utest 

detectlocal:
	cd scripts; bash detect-local-gateway-repo.sh

detectmisconfig:
	cd scripts; bash detect-pkg-misconfig.sh

useremote:
	cd scripts; bash use-remote-repos.sh

uselocal:
	echo "replace github.com/ConsenSys/fc-retrieval/common => ../fc-retrieval/common" >> go.mod
	go mod tidy

utest:
	go test ./...

coverage:
	bash ./scripts/coverage.sh $(COV)

# Alays assume these targets are out of date.
.PHONY: clean itest utest build release push detectmisconfig detectlocal stop coverage

