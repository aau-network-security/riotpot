# Makefile
APPNAME=riotpot
CURRENT_DIR=`pwd`
PACKAGE_DIRS=`go list -e ./... | egrep -v "binary_output_dir|.git|mocks"`
DOCKERFILE=docker/Dockerfile.documentation

# docker cmd below
.PHONY:  docker-build docker-run
docker-build:
	docker build -f $(DOCKERFILE) . -t $(APPNAME)/v1
docker-run: docker-build
	docker run -p 6060:6060 -it $(APPNAME)/v1