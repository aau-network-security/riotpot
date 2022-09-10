FROM golang:1.18 as builder
LABEL maintainer="Shreyas Srinivasa <shsr@es.aau.dk>"

# Set the working directory to golang working space
WORKDIR /app

# Copy the dependencies into the image
COPY go.mod ./
COPY go.sum ./

# download all the dependencies
RUN go mod download

# Copy everything into the image
# Copy only the app files in the image
COPY api api/
COPY cmd cmd/
COPY internal internal/
COPY pkg pkg/
COPY tools tools/

# Copy static files
COPY statik/ statik/

ADD Makefile .
RUN make riotpot-builder

# Run the command from the Makefile to build all the plugins
# and build the project
# -- Comment this line when on development if you know you have a ready to go version built --
# Disclaimer: if you comment this line, be 100% sure that the binary can be run on linux
FROM gcr.io/distroless/base-debian10

WORKDIR /

# Copy the dependencies into the image
# TODO: Copy the plugins
COPY --from=builder /riotpot /riotpot
EXPOSE 2022

# Run RIoTPot
ENTRYPOINT [ "/riotpot" ]
