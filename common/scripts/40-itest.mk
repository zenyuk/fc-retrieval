.PHONY: test itest

quick-itest test: itest-poc1

itest-poc1:
	$(MAKE) LAST_GATEWAY_NO=0 _mk_itest-poc1 clean

itest-poc2js:
	$(MAKE) LAST_GATEWAY_NO=1 _mk_itest-poc2js clean

_mk_itest-% \
:
	set -e; \
	make docker-clean docker-restart; \
	DIR=$(subst _mk_itest-,,$@); \
	cd itest; \
	. ./env-LOTUS_KEYS.out; \
	cd pkg/$$DIR; \
	echo " \\e[01;32m \\n#run: $@\\e[m"; \
	go test -count=1 -p 1 -v -failfast \
	-covermode count -coverprofile cover.out \
	-coverpkg github.com/ConsenSys/fc-retrieval/common/...; \
	free -h

