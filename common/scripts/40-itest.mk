.PHONY: test itest

# client-gateway  client-init  lotus  lotus-full-node  poc1  poc2_dht_offer  poc2_dht_offer_ack  poc2_dht_offer_new_gateway  poc2_group_offer  poc2js  poc2_new_gateway  provider-admin  util
#                                     lotus-full-node                                      :x    poc2_dht_offer_new_gateway  poc2_group_offer  poc2js  poc2_new_gateway                  util
#
quick-itest test:
	@echo " \\e[01;32m \\n#run: $@\\e[m"
	$(MAKE) LAST_GATEWAY_NO=0 itest-client-gateway
	$(MAKE) LAST_GATEWAY_NO=0 itest-client-init
	$(MAKE) LAST_GATEWAY_NO=0 itest-provider-admin
	$(MAKE) LAST_GATEWAY_NO=0 itest-poc1
	true $(MAKE) LAST_GATEWAY_NO=0 itest-lotus
	true $(MAKE) LAST_GATEWAY_NO=0 itest-lotus-full-node
	true $(MAKE) LAST_GATEWAY_NO=0



itest-% \
:
	@echo " \\e[01;32m \\n#run: $@\\e[m"
	set -e; \
	make docker-clean docker-restart; \
	DIR=$(subst itest-,,$@); \
	cd itest; \
	. ./env-LOTUS_KEYS.out; \
	cd pkg/$$DIR; \
	echo " \\e[01;32m \\n#run: $@\\e[m"; \
	go test -count=1 -p 1 -v -failfast \
	-covermode count -coverprofile cover.out \
	-coverpkg github.com/ConsenSys/fc-retrieval/common/...; \
	free -h

