REGISTRY?=
VERSION?=v1

start:
	docker-compose up

stop:
	docker-compose stop

build:
	docker build -f Dockerfile.dev -t ${REGISTRY}fc-retrieval-register:${VERSION} .

build-dev:
	go build -v cmd/filecoin-retrieval-register-server/main.go