# Makefile
APPNAME=riotpot
DEPLOY=deployments/
DOCKER=build/docker/
PLUGINS_DIR=pkg/plugin

# docker cmd below
.PHONY:  riotpot-install docker-build-doc riotpot-doc riotpot-up riotpot-build riotpot-build-plugins riotpot-builder
docker-build-doc:
	docker build -f $(DOCKER)Dockerfile.documentation . -t $(APPNAME)/v1
riotpot-doc: docker-build-doc
	docker run -p 6060:6060 -it $(APPNAME)/v1
riotpot-up:
	docker-compose -p riotpot -f ${DEPLOY}docker-compose.yml up -d --build
riotpot-down:
	docker-compose -p riotpot -f ${DEPLOY}docker-compose.yml down -v
riotpot-all:
	riotpot-doc
	riotpot-up
riotpot-build:
	go build -o ./riotpot ./cmd/riotpot/.
riotpot-build-plugins: $(PLUGINS_DIR)/*
	for folder in $^ ; do \
		go build -buildmode=plugin -o $${folder}/plugin.so $${folder}/*.go; \
	done
riotpot-builder: \
	riotpot-build \
	riotpot-build-plugins
riotpot-ui:
	@cd ui && serve -s build
	