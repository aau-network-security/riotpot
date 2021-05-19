# Makefile
APPNAME=riotpot
CURRENT_DIR=`pwd`
PACKAGE_DIRS=`go list -e ./... | egrep -v "binary_output_dir|.git|mocks"`
DEPLOY=deployments/
DOCKER=build/docker/
PLUGINS_DIR=pkg/plugin

# docker cmd below
.PHONY:  docker-build-doc riotpot-doc riotpot-up riotpot-prod-up riotpot-prod-down riotpot-build riotpot-build-plugins riotpot-builder
docker-build-doc:
	docker build -f $(DOCKER)Dockerfile.documentation . -t $(APPNAME)/v1
riotpot-doc: docker-build-doc
	docker run -p 6060:6060 -it $(APPNAME)/v1
riotpot-up:
	docker-compose -p riotpot -f ${DEPLOY}docker-compose.yml up -d --build
riotpot-down:
	docker-compose -p riotpot -f ${DEPLOY}docker-compose.yml down -v
riotpot-prod-up:
	docker-compose -p riotpot -f ${DEPLOY}docker-compose.prod.yml up -d --build
riotpot--prod-down:
	docker-compose -p riotpot -f ${DEPLOY}docker-compose.prod.yml down -v
riotpot-all:
	riotpot-doc
	riotpot-up
riotpot-build:
	go build -o riotpot cmd/riotpot/main.go;
riotpot-build-plugins: $(PLUGINS_DIR)/*
	for folder in $^ ; do \
		go build -buildmode=plugin -o $${folder}/plugin.so $${folder}/*.go; \
	done
riotpot-builder: \
	riotpot-build \
	riotpot-build-plugins