FROM node:19.4.0-alpine as ui

# set working directory
WORKDIR /app

# add `/app/node_modules/.bin` to $PATH
ENV PATH /app/node_modules/.bin:$PATH

# install app dependencies
# add app
COPY ui/ .
RUN npm install --silent

# Creates the build
RUN npm run build --omit=dev

FROM golang:1.18 as builder

# Set the working directory to golang working space
WORKDIR /riotpot

# Copy the build and `embed.go` files from the ui folder
COPY --from=ui app/build/ ui/build/
COPY --from=ui app/embed.go ui/

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

RUN make build-all

# Run the command from the Makefile to build all the plugins
# and build the project
# -- Comment this line when on development if you know you have a ready to go version built --
# Disclaimer: if you comment this line, be 100% sure that the binary can be run on linux
FROM gcr.io/distroless/base-debian10

ENV GIN_MODE=release

WORKDIR /riotpot

# Copy the dependencies into the image
COPY --from=builder /riotpot/bin/ ./

# API, required for the UI.
EXPOSE 2022

CMD ["./riotpot"]