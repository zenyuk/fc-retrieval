VERSION?=0.0.1

dev:
	docker-compose -f docker-compose.dev.yml up

# stop:
# 	docker-compose stop

# build:
# 	docker-compose build

build-dev:
	go build -v cmd/filecoin-retrieval-register-server/main.go