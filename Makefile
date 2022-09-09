# Makefile
APPNAME=riotpot
CURRENT_DIR=`pwd`
PACKAGE_DIRS=`go list -e ./... | egrep -v "binary_output_dir|.git|mocks"`
DEPLOY=deployments/
DOCKER=build/docker/
PLUGINS_DIR=pkg/plugin
LCL_GO_PATH=$(shell echo `go env GOPATH`)
GO_BIN_DIR=$(LCL_GO_PATH)/bin
LOCAL_BUILD_ENABLED=False

# docker cmd below
.PHONY:  docker-build-doc riotpot-doc riotpot-up riotpot-prod-up riotpot-prod-down riotpot-build riotpot-build-plugins riotpot-builder riotpot-build-local set_container_build set_local_build go_install
prepare_configurations:
	cp configs/samples/configuration-template.yml configs/samples/configuration.yml
set_local_build:
	@sed -i -e 's/local_build_on: 0/local_build_on: 1/g' configs/samples/configuration.yml
set_container_build:
	@sed -i -e 's/local_build_on: 1/local_build_on: 0/g' configs/samples/configuration.yml
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
	go build -o riotpot ./cmd/riotpot/.
riotpot-build-plugins: $(PLUGINS_DIR)/*
	for folder in $^ ; do \
		go build -buildmode=plugin -o $${folder}/plugin.so $${folder}/*.go; \
	done
riotpot-build-local-plugin: $(PLUGINS_DIR)/*
	for folder in $^ ; do \
		go build -buildmode=plugin -o ${GO_BIN_DIR}/$${folder}/plugin.so $${folder}/*.go; \
	done

riotpot-builder: \
	set_container_build \
	riotpot-build \
	riotpot-build-plugins
riotpot-build-local: \
	riotpot-build \
	riotpot-build-local-plugin \
	go install ./riotpot
