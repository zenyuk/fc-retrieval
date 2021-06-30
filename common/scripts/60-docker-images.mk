
# make micro docker images
micro-images: register/main provider/main gateway/main
	set -xe; \
        for d in register provider gateway; \
        do \
		cd $$d/micro-docker; \
                docker rmi consensys/fc-retrieval/$$d 2>/dev/null || true; \
		cp -uv ../main .; \
		docker build --no-cache -t consensys/fc-retrieval/$$d . ; \
		cd ../../; \
        done;
