# Copyright (C) 2020 ConsenSys Software Inc.

# Build the Filecoin Retrieval Gateway

default: clean utest

# builds a docker image that builds the app and packages it into a minimal docker image

utest:
	go test ./...

clean:
#nothing to do

useremote:
	cd scripts; bash use-remote-repos.sh

.PHONY: default clean utest

