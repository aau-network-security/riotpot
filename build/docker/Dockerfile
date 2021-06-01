# Start from the latest golang base image
FROM golang:1.16
LABEL maintainer="Shreyas Srinivasa <shsr@es.aau.dk>"

ENV GO111MODULE=on
ENV APP_NAME=riotpot
ENV WORKING_DIR=$GOPATH/src/github.com/aau/$APP_NAME

# Set the working directory to golang working space
WORKDIR $WORKING_DIR

# Copy the dependencies into the image
COPY go.mod .
COPY go.sum .

# download all the dependencies
RUN go mod download

# install netcat so we can ping the db until it is ready
RUN apt-get update && apt-get install -y netcat

# Copy everything into the image
# TODO: this can be optimized to copy just the necessary files!
COPY . .

# Run the command from the Makefile to build all the plugins
# and build the project
# -- Comment this line when on development if you know you have a ready to go version built --
# Disclaimer: if you comment this line, be 100% sure that the binary can be run on linux
RUN make riotpot-builder

# give permissions to the entrypoint to run the file
RUN chmod +x build/docker/entrypoint.sh

# Run RIoTPot
ENTRYPOINT [ "./build/docker/entrypoint.sh" ]
