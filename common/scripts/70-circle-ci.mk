#
# make ci: it will run steps, which runs on circle CI before normal itest
#  

ci-1 : deps build-servers micro-images coverage

ci-2 ci-lotus:
	set -e; \
	docker image ls | egrep '^REPOSITORY|consensys'; \
	echo; \
	cd itest; \
	make lotusbase lotusdaemon lotusfullnode

ci-3 : build-servers-test

ci-4 ci-poc1 :
	$(MAKE) LAST_GATEWAY_NO=0 itest-poc1 clean


ciall: ci-1 ci-2 ci-3 ci-4

# skip lotus and any itest
ci: ci-1 ci-3
