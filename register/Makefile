VERSION?=dev
IMAGE?=consensys/fc-retrieval-register

default: clean build tag

# User `make dev arg=--build` to rebuild
dev:
	make -s env-soft
	GO_MOD=go.mod docker-compose -f docker-compose.dev.yml up $(arg)

dev-local:
	make -s env-soft
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

env-soft:
	make env-compare
	if [ -s .env ]; then \
		echo ".env already exists"; \
	else \
		echo ".env does not exists, create a .env file"; \
		cp .env.example .env; \
	fi

env-compare:
	@cmp -s .env .env.example; \
	RETVAL=$$?; \
	if [ $$RETVAL -eq 0 ]; then \
		echo ".env and .env.example are the same"; \
	else \
		echo ".env and .env.example are NOT the same"; \
	fi