# Copyright (C) 2020 ConsenSys Software Inc.

# Build the Filecoin Retrieval Provider Admin
COV?=80

# This target (the first target in the build file) is the one that is executed if no 
# command line args are specified.
release: clean utest build

# builds a docker image that builds the app and packages it into a minimal docker image
build:
	go build ./...

uselocal:
	cd scripts; bash use-local-repos.sh

useremote:
	cd scripts; bash use-remote-repos.sh

detectlocal:
	cd scripts; bash detect-local-gateway-repo.sh

detectmisconfig:
	cd scripts; bash detect-gateway-misconfig.sh

utest:
	go test ./...

coverage:
	bash ./scripts/coverage.sh $(COV)


# remove previous images and containers
clean:

# Alays assume these targets are out of date.
.PHONY: clean itest utest coverage build release push detectmisconfig detectlocal

