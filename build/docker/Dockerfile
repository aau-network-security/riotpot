FROM golang:1.18 as builder
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

# Copy everything into the image
# Copy only the app files in the image
COPY internal internal/
COPY tools tools/
COPY pkg pkg/
COPY cmd cmd/
COPY configs configs/
COPY build/docker/entrypoint.sh ./

# Run the command from the Makefile to build all the plugins
# and build the project
# -- Comment this line when on development if you know you have a ready to go version built --
# Disclaimer: if you comment this line, be 100% sure that the binary can be run on linux
COPY Makefile .
RUN make builder

FROM golang:1.18 as stage
WORKDIR /app

# TODO: Copy the plugins
COPY --from=stage cmd/riotpot ./
USER 1000

# Run RIoTPot
ENTRYPOINT [ "." ]
