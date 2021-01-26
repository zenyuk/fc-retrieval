PORT ?= 8081
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

dev:
	docker-compose -f docker-compose.dev.yml up

stop:
	docker-compose stop

utest:
	go test ./...
