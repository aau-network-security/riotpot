
<div align="center">
  <img src="AAUgraphics/aau_logo.png" height="100">
  <p align="center">
    <h2 align="center">RiotPot</h2>
  </p>
  <p align="center">
    <!-- Workflow status -->
    <a href="https://github.com/aau-network-security/riotpot/actions"><img alt="GitHub Actions status" src="https://github.com/aau-network-security/riotpot/workflows/cyber/badge.svg"></a>
    <a href="https://GitHub.com/aau-network-security/riotpot/releases/"><img alt="GitHub Actions status" src="https://img.shields.io/github/release/aau-network-security/riotpot"></a>
    <a href="https://github.com/aau-network-security/riotpot/"><img alt="GitHub Actions status" src="https://img.shields.io/github/go-mod/go-version/aau-network-security/riotpot"></a>
  </p>
</div>

___

- [1. Description](#1-description)
- [2. Installation](#2-installation)
  - [2.1 Docker](#21-docker)
    - [2.1.1 Docker Hub Image](#211-docker-hub-image)
  - [2.2 Local](#22-local)

## 1. Description

Riotpot is an interoperable medium interaction honeypot, primarly focused on the emulation IoT and OT protocols, although, it is also capable of emulating other services.

This services are loaded in the honeypot in the form of plugins, making Riotpot a modular, and very transportable honeypot. The services are loaded at runtime, meaning that the weight of the honeypot will vary on premisses, and the services loaded e.g. HTTP, will only be used when required. As consequence, we higly recommend building your own binary customised to your own needs. Refer to the following section, Installation, for more information.

## 2. Installation

Although one can donwload the binary containing a set of default running emulators, this guide is mainly focused to those looking for a customized experience.

We thrive on the idea of making RiotPot higly transportable, therefore, in this section one can find multiple methods of installation for diverse environments that fit a broad list of requirements and constrains.

We highly recommend running RiotPot in a virtualized self-contained network using Docker, for which we included configuration files that run the honeypot both as a closed environment for testing and playing around (similar to a testbed environment), and a production environment that runs on the host machine exposing services to the host main network that can run from the local clone of this repository. 

> **NOTE:** The production image can be pulled from Docker Hub. If you choose this method you may directly jump to [2.1 Docker](#21-docker).
<!-- Not implemented yet, include section for the Docker Image pulling and building -->

RiotPot is written in Golang, therefore, you will need to have [go](https://golang.org) installed first if you plan to make any changes, otherwise you can skip steps 1 and 2 if you rather not installing go. 

Regardless, you will need to copy riotpot to local:

```bash
# 1. Make the folder in where the repository will be stored.
mkdir -p $GOPATH/src/github.com

# 2. Navigate to the folder in where you store your repositories
mkdir -p $GOPATH/src/github.com

# 3. Clone the repository
git clone git@github.com:aau-network-security/riotpot.git

# 4. Navigate to the newly created folder with the repository
cd riotpot
```

### 2.1 Docker

We assume you have basic knowledge about the Docker ecosystem, otherwise please refer first to the Docker documentation [here](https://docs.docker.com).

At the root folder of RiotPot there are two **docker-compose** files:
```bash
cd ~/riotpot | ls -al
...
-rw-r--r--  docker-compose.dev.yml
-rw-r--r--  docker-compose.yml
...
```

This files correspond to the respective software development environments: *development* and *production*.

**Development.**
The first one, *docker-compose.dev.yml* builds the project in a private virtual network in which there are two hosts: *riotpot*, and *attacker*. Both can only communicate with the each other. Use this setup for **development** and **testing locally** by typing in your a terminal:

```bash
docker-compose -f docker-compose.dev.yml up -d --build
```

**Production.**
The second, *docker-compose.yml* is a production ready environment which is only constructed using one host, **riotpot**. Use the following command to put build the container: 

```bash
docker-compose up -d --build
```

Once you are done with the honeypot, you can put down the containers using the *down* command. 

```bash
docker-compose down -v
```

> **NOTE:** Using the *-v* tag will remove all the mounted volumes, i.e. the database used by riotpot to store information and the volumes mounted to store logs and binaries collected by the honeypot.

#### 2.1.1 Docker Hub Image

Build the latest release of RiotPot directly from the image provided in the Docker Hub:

```bash
# Grab and run the latest release of the riotpot consummer immage
# dettached from the console with -d.
docker run -d riotpot-docker:latest
```

### 2.2 Local

To build your own binary from source, navigate to the folder where you have stored the repository and use the go CLI to generate it and store it in the *./bin/* folder:

```bash
# build the binary in the ./bin folder
go build -o bin/
```

Additionally, you could also install the application in the system:
```bash
# installs riotpot at $GOPATH/bin
go install
```

Run the binary as any other application:
```bash
./path/to/riotpot
```
