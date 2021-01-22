PORT ?= 8080
REGISTRY?=
VERSION?=v1

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

start-dev:
	go run cmd/provider/main.go --host 0.0.0.0 --port $(PORT)

stop:
	docker-compose stop

utest:
	go test ./...
