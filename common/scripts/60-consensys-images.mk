
image-%:
	set -e; \
	DIR=$(subst image-,,$@); \
	docker build -t consensys/fc-retrieval/$$DIR:dev -f $$DIR/Dockerfile . ; \
	cd $$DIR; \
	make tag
