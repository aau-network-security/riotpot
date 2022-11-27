# Makefile
APPNAME=riotpot
DEPLOY=deployments/
DOCKER=build/docker/
PLUGINS_DIR=pkg/plugin

# docker cmd below
.PHONY:  docker-build-doc docker-doc-up up down up-all build build-plugins build-all ui
docker-build-doc:
	docker build -f $(DOCKER)Dockerfile.documentation . -t $(APPNAME)/v1
docker-doc-up: docker-build-doc
	docker run -p 6060:6060 -it $(APPNAME)/v1
up:
	docker-compose -p riotpot -f ${DOCKER}docker-compose.yaml up -d --build
down:
	docker-compose -p riotpot -f ${DOCKER}docker-compose.yaml down -v
up-all:
	riotpot-doc
	riotpot-up
build:
	go build -gcflags='all=-N -l' -o ./bin ./cmd/riotpot/.
build-plugins: $(PLUGINS_DIR)/*
	for folder in $^ ; do \
		result=$${folder%%+(/)}; \
		result=$${result##*/}; \
		result=$${result:-/}; \
		go build -buildmode=plugin --mod=vendor -gcflags='all=-N -l' -o bin/plugins/$${result}.so $${folder}/*.go; \
	done
build-all: \
	build \
	build-plugins
ui:
	@cd ui && serve -s build