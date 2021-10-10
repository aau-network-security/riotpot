# Makefile
APPNAME=riotpot
CURRENT_DIR=`pwd`
PACKAGE_DIRS=`go list -e ./... | egrep -v "binary_output_dir|.git|mocks"`
DEPLOY=deployments/
DOCKER=build/docker/
PLUGINS_DIR=pkg/plugin

# docker cmd below
.PHONY:  docker-build-doc doc up prod-up prod-down build plugins builder tor-up tor-down
docker-build-doc:
	docker build -f $(DOCKER)Dockerfile.documentation . -t $(APPNAME)/v1
doc: docker-build-doc
	docker run -p 6060:6060 -it $(APPNAME)/v1
up:
	docker-compose -p riotpot -f ${DEPLOY}docker-compose.yml up -d --build
down:
	docker-compose -p riotpot -f ${DEPLOY}docker-compose.yml down -v
prod-up:
	docker-compose -p riotpot -f ${DEPLOY}docker-compose.prod.yml up -d --build
prod-down:
	docker-compose -p riotpot -f ${DEPLOY}docker-compose.prod.yml down -v
tor-up:
	docker-compose -p riotpot -f ${DEPLOY}docker-compose.tor.yml up -d --build
tor-down:
	docker-compose -p riotpot -f ${DEPLOY}docker-compose.tor.yml down -v
riotpot-all:
	doc
	up
build:
	go build -o riotpot cmd/riotpot/main.go;
plugins: $(PLUGINS_DIR)/*
	for folder in $^ ; do \
		go build -buildmode=plugin -o $${folder}/plugin.so $${folder}/*.go; \
	done
builder: \
	build \
	plugins