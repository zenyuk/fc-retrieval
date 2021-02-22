VERSION?=dev
IMAGE?=consensys/fc-retrieval-register

default: clean build tag

# User `make dev arg=--build` to rebuild
dev:
	GO_MOD=go.mod docker-compose -f docker-compose.dev.yml up $(arg)

dev-local:
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
