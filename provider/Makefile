PORT ?= 9030
REGISTRY?=consensys/
VERSION?=dev

.PHONY: build build-dev start start-dev stop

build:
	docker build -t ${REGISTRY}fc-retrieval-provider:${VERSION} .

build-dev:
	go build -v cmd/provider/main.go

clean:
	docker rm -f ${REGISTRY}fc-retrieval-provider:${VERSION} 2> /dev/null || true
	docker rmi -f "${REGISTRY}fc-retrieval-provider:${VERSION}" || true

start:
	docker-compose up -d

# User `make dev arg=--build` to rebuild
dev:
	docker-compose -f docker-compose.dev.yml up $(arg)

stop:
	docker-compose stop

utest:
	go test ./...

uselocal:
	cd scripts; bash use-local-repos.sh

useremote:
	cd scripts; bash use-remote-repos.sh