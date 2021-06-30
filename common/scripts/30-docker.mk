
.PHONY:     lotus-full-node-image docker-clean docker-restart
docker-all: lotus-full-node-image docker-clean docker-restart

.PHONY:         docker-network docker-lotus-full-node docker-redis docker-register docker-provider-test docker-gateway-test docker-hosts docker-itest-env
docker-restart: docker-network docker-lotus-full-node docker-redis docker-register docker-provider-test docker-gateway-test docker-hosts docker-itest-env

lotus-base-image:
	@echo " \\e[01;32m \\n#run: $@\\e[m"
	set -e; \
	cd itest; \
	test -n "$$(docker image ls consensys/lotus-base -q)" || make lotusbase

lotus-daemon-image: lotus-base-image
	@echo " \\e[01;32m \\n#run: $@\\e[m"
	set -e; \
	cd itest; \
	test -n "$$(docker image ls consensys/lotus-daemon -q)" || make lotusdaemon

lotus-full-node-image: lotus-daemon-image
	@echo " \\e[01;32m \\n#run: $@\\e[m"
	set -e; \
	cd itest; \
	test -n "$$(docker image ls consensys/lotus-full-node -q)" || make lotusfullnode

.PHONY: clean
clean docker-clean:
	@echo " \\e[01;32m \\n#run: $@\\e[m"
	(docker ps -q    | xargs docker stop 2>/dev/null) || true
	(docker ps -q -a | xargs docker rm   2>/dev/null) || true

docker-network:
	@echo " \\e[01;32m \\n#run: $@\\e[m"
	docker network create dokbr01 || true

docker-lotus-daemon:
	@echo " \\e[01;32m \\n#run: $@\\e[m"
	docker run -d --network dokbr01 --name lotus-daemon consensys/lotus-daemon

docker-lotus-full-node:
	@echo " \\e[01;32m \\n#run: $@\\e[m"
	docker run -d --network dokbr01 --name lotus-full-node consensys/lotus-full-node

docker-redis:
	@echo " \\e[01;32m \\n#run: $@\\e[m"
	docker run --name redis --network dokbr01 \
		-e ALLOW_EMPTY_PASSWORD=yes \
		redis:alpine redis-server --requirepass xxxx >redis.out 2>&1 &
	
docker-register:
	@echo " \\e[01;32m \\n#run: $@\\e[m"
	cd register; \
	cp -u .env.example .env; \
	docker run -d \
		--env-file .env \
		--network dokbr01 \
		-v $$(pwd):/app \
		--name register \
		-w /app \
		-e DOCKER_NAME=$$dname \
		-e REDIS_PASSWORD=xxxx \
		redhat/ubi8-micro sh -c "./main --host=0.0.0.0 --port=9020 >register.out 2>&1"

docker-provider:
	@echo " \\e[01;32m \\n#run: $@\\e[m"
	set -ep; \
	cd provider; \
	realpath main; \
	cp -u .env.example .env; \
	for I in "" `printf "2\n" $(LAST_GATEWAY_NO) | sort -n | head -1 | xargs seq 0| sed 's/^/-/'`; \
	do \
	dname=provider$$I; \
	docker run -d \
		--env-file .env \
		--network dokbr01 \
		-v $$(pwd):/app \
		--name $$dname \
		-w /app \
		-e DOCKER_NAME=$$dname \
		-e REDIS_PASSWORD=xxxx \
		redhat/ubi8-micro sh -c "./main >$$dname.out 2>&1"; \
	done

docker-gateway:
	@echo " \\e[01;32m \\n#run: $@\\e[m"
	cd gateway; \
	cp -u .env.example .env; \
	for I in "" `seq 0 32| sed 's/^/-/'`; \
	do \
	dname=gateway$$I; \
	docker run -d \
		--env-file .env \
		--network dokbr01 \
		-v $$(pwd):/app \
		--name $$dname \
		-w /app \
		-e DOCKER_NAME=$$dname \
		-e REDIS_PASSWORD=xxxx \
		redhat/ubi8-micro sh -c "./main >$$dname.out 2>&1"; \
	done

docker-provider-test:
	@echo " \\e[01;32m \\n#run: $@\\e[m"
	cd provider; \
	realpath main; \
	cp -u .env.example .env; \
	for I in "" `printf "2\n" $(LAST_GATEWAY_NO) | sort -n | head -1 | xargs seq 0| sed 's/^/-/'`; \
	do \
	dname=provider$$I; \
	docker run -d \
		--env-file .env \
		--network dokbr01 \
		-v $$(pwd):/app \
		--name $$dname \
		-w /app \
		-e DOCKER_NAME=$$dname \
		-e REDIS_PASSWORD=xxxx \
		redhat/ubi8-micro sh -c "./main.test -test.coverprofile $$dname.cov >$$dname.out 2>&1"; \
	done

LAST_GATEWAY_NO?=32
docker-gateway-test:
	@echo " \\e[01;32m \\n#run: $@\\e[m"
	cd gateway; \
	cp -u .env.example .env; \
	for I in "" `seq 0 $(LAST_GATEWAY_NO)| sed 's/^/-/'`; \
	do \
	dname=gateway$$I; \
	docker run -d \
		--env-file .env \
		--network dokbr01 \
		-v $$(pwd):/app \
		--name $$dname \
		-w /app \
		-e DOCKER_NAME=$$dname \
		-e REDIS_PASSWORD=xxxx \
		redhat/ubi8-micro sh -c "./main.test -test.coverprofile $$dname.cov >$$dname.out 2>&1"; \
	done

docker-hosts:
	@echo " \\e[01;32m \\n#run: $@\\e[m"
	cd itest; \
	docker ps -q \
        | xargs -n 1 docker inspect --format \
        '{{ .Name }} {{range .NetworkSettings.Networks}} {{.IPAddress}}{{end}}' \
        | sed 's#^/##' | sed 's/$$/.nip.io/' >./hosts.out

docker-itest-env:
	@echo " \\e[01;32m \\n#run: $@\\e[m"
	cd itest; \
	docker exec  lotus-full-node bash -c \
	  'echo export LOTUS_TOKEN=$$(cat ~/.lotus/token); echo; echo export SUPER_ACCT=$$(/app/lotus/lotus wallet default)' > \
	  env-LOTUS_KEYS.out; \
	echo export ITEST_CALLING_FROM_CONTAINER=yes >> env-LOTUS_KEYS.out; \
	echo export LOG_LEVEL=debug >> env-LOTUS_KEYS.out; \
	echo export LOG_TARGET=STDOUT >> env-LOTUS_KEYS.out; \
	echo export LOG_SERVICE_NAME=itest >> env-LOTUS_KEYS.out; \
	echo export HOSTALIASES=`realpath hosts.out` >> env-LOTUS_KEYS.out;
