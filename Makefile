# Makefile
APPNAME=riotpot
CURRENT_DIR=`pwd`
PACKAGE_DIRS=`go list -e ./... | egrep -v "binary_output_dir|.git|mocks"`
DOCKERFILE_DOC=docker/Dockerfile.documentation

# docker cmd below
.PHONY:  docker-build-doc riotpot-doc riotpot-up riotpot-dev riotpot-dev-all
docker-build-doc:
	docker build -f $(DOCKERFILE_DOC) . -t $(APPNAME)/v1
riotpot-doc: docker-build-doc
	docker run -p 6060:6060 -it $(APPNAME)/v1
riotpot-up:
	docker-compose up -d --build
riotpot-dev:
	docker-compose -f docker.compose.dev.yml up -d --build
riotpot-dev-all:
	riotpot-doc
	riotpot-dev