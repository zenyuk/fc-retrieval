
build: timed-deps coverage build-servers build-test-binary build-servers-test

deps:
	@echo " \\e[01;32m \\n#run: $@\\e[m"
	@set -e; \
	for DIR in $(wildcard */go.mod); \
	do \
	  cd $$(dirname $$DIR); \
	  echo running go mod tidy in dir:; \
	  pwd; \
	  go mod tidy; \
	  cd ..; \
	done; \
	echo; \
	echo go.sum below should be updated: ; \
	for DIR in $(wildcard */go.mod); \
	do \
	  cd $$(dirname $$DIR); \
	  git -P diff -U0 --exit-code go.sum || pwd; \
	  cd ..; \
	done; \
	echo; \
	echo go.mod below should be updated: ; \
	for DIR in $(wildcard */go.mod); \
	do \
	  cd $$(dirname $$DIR); \
	  git -P diff -U0 --exit-code go.mod || pwd; \
	  cd ..; \
	done;

timed-deps:
	@echo " \\e[01;32m \\n#run: $@\\e[m"
	/usr/bin/time -vv make deps

coverage:
	@echo " \\e[01;32m \\n#run: $@\\e[m"
	set -e; \
	cd common; \
	make coverage COV=62; \
	make utest

build-servers:
	@echo " \\e[01;32m \\n#run: $@\\e[m"
	set -e; \
	for D in gateway provider register; do cd $$D; go build ./cmd/*/main.go; cd ..; done

build-test-binary:
	@echo " \\e[01;32m \\n#run: $@\\e[m"
	set -e; \
	cd itest; \
	go list -f '{{.Dir}}' ./... | while read D; do cd $$D; go test -c -covermode count $$D; done

build-servers-test:
	@echo " \\e[01;32m \\n#run: $@\\e[m"
	set -e; \
	for D in gateway provider register; do cd $$D; \
	cd cmd/$$D*; \
	cd ../..; \
	go test -c ./cmd/*/. -o main.test \
		-covermode count \
		-coverpkg github.com/ConsenSys/fc-retrieval/common/... ; \
	cd ..; \
	done
