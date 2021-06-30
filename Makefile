
default:

ci-part1 : deps build-servers micro-images coverage
ci-part2 :
ci-part3 : build-servers-test
ci-part4 :
	$(MAKE) LAST_GATEWAY_NO=0 itest-poc1 clean


ci: ci-part1 ci-part2 ci-part3 ci-part4

-include ./common/scripts/*.mk
-include ../local*.mk
