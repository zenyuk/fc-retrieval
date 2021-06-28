VERSION?=dev
IMAGE?=consensys/fc-retrieval-client-js

.PHONY: build clean tag

default: clean build tag

build:
	docker build -t $(IMAGE):$(VERSION) .

clean:
	docker rm -f $(IMAGE):$(VERSION) 2> /dev/null || true
	docker rmi -f "$(IMAGE):$(VERSION)" || true

tag:
	cd scripts; bash tag.sh $(VERSION) $(IMAGE):$(VERSION)
