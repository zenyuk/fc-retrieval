

lotus-images:
	@echo " \\e[01;32m \\n#run: $@\\e[m"
	set -e; \
	docker image ls | egrep '^REPOSITORY|consensys'; \
	echo; \
	cd itest; \
	make lotusbase lotusdaemon lotusfullnode

# make micro docker images
micro-images: register/main provider/main gateway/main
	@echo " \\e[01;32m \\n#run: $@\\e[m"
	set -xe; \
	for d in register provider gateway; \
	do \
		cd $$d/micro-docker; \
		cp -uv ../main .; \
		docker rmi consensys/fc-retrieval/$$d 2>/dev/null || true; \
		docker build --no-cache -t consensys/fc-retrieval/$$d . ; \
		cd ../../; \
	done;
