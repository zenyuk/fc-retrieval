PORT ?= 9030
VERSION?=dev
IMAGE?=consensys/fc-retrieval-provider
COV?=80

.PHONY: build build-dev start start-dev stop utest coverage

default: clean build tag

build:
	docker build -t $(IMAGE):$(VERSION) .

build-dev:
	go build -v cmd/provider/main.go

build-local:
	cat go.mod >> temp
	echo "replace github.com/ConsenSys/fc-retrieval-common => ./local/fc-retrieval-common" >> go.mod
	rm -rf ./local/
	mkdir -p ./local/fc-retrieval-common/pkg
	cp -r ../fc-retrieval-common/pkg/ ./local/fc-retrieval-common/pkg/
	cp ../fc-retrieval-common/go.mod ./local/fc-retrieval-common/go.mod
	docker build -t $(IMAGE):$(VERSION) .
	rm -rf ./local/
	mv temp go.mod
	
push:
	cd scripts; bash push.sh $(VERSION) $(IMAGE):$(VERSION)

tag:
	cd scripts; bash tag.sh $(VERSION) $(IMAGE):$(VERSION)

clean:
	docker rm -f $(IMAGE):$(VERSION) 2> /dev/null || true
	docker rmi -f "$(IMAGE):$(VERSION)" || true

start:
	docker-compose up -d

# User `make dev arg=--build` to rebuild
dev:
	./scripts/make-env-file/make-env-file -source=.env.example -dest=.env
	docker-compose -f docker-compose.dev.yml up $(arg)

stop:
	docker-compose stop

utest:
	go test ./...

coverage:
	bash ./scripts/coverage.sh $(COV)

uselocal:
	cd scripts; bash use-local-repos.sh

useremote:
	cd scripts; bash use-remote-repos.sh
