VERSION?=dev
IMAGE?=consensys/fc-retrieval-register
COV?=80

default: clean build tag

# User `make dev arg=--build` to rebuild
dev:
	./scripts/make-env-file/make-env-file -source=.env.example -dest=.env
	GO_MOD=go.mod docker-compose -f docker-compose.dev.yml up $(arg)

dev-local:
	./scripts/make-env-file/make-env-file -source=.env.example -dest=.env
	GO_MOD=go.local.mod docker-compose -f docker-compose.dev.yml up $(arg)

# stop:
# 	docker-compose stop

build:
	docker build -f Dockerfile -t ${IMAGE}:${VERSION} .

build-local:
	docker build -f Dockerfile.dev -t ${IMAGE}:${VERSION} .

build-dev:
	go build -v cmd/register-server/main.go

push:
	cd scripts; bash push.sh ${VERSION} ${IMAGE}:${VERSION}

tag:
	cd scripts; bash tag.sh ${VERSION} ${IMAGE}:${VERSION}

uselocal:
	cd scripts; bash use-local-repos.sh

useremote:
	cd scripts; bash use-remote-repos.sh

# Have a clean target to match the other repos. This will be called by the integration test 
# system when building the register
clean:
	echo Does nothing

utest:
	go test ./...

coverage:
	bash ./scripts/coverage.sh $(COV)

# Alays assume these targets are out of date.
.PHONY: dev dev-local build build-local build-dev push tag uselocal useremote clean utest coverage
