
<div align="center">
  <img src="assets/aau_logo.png" height="100">
  <p align="center">
    <h2 align="center">RiotPot</h2>
  </p>
  <p align="center">
    <!-- Workflow status -->
    <a href="https://github.com/aau-network-security/riotpot/actions"><img alt="GitHub Actions status" src="https://github.com/aau-network-security/riotpot/workflows/cyber/badge.svg"></a>
    <a href="https://goreportcard.com/badge/github.com/aau-network-security/riotpot"><img src="https://goreportcard.com/badge/github.com/aau-network-security/riotpot?style=flat-square"></a>
    <a href="https://pkg.go.dev/riopot"><img src="https://pkg.go.dev/badge/riopot.svg"></a>
    <a href=""><img src="https://img.shields.io/github/release/riotpot/project-layout.svg?style=flat-square">
  </p>
</div>
___

- [1. Description](#1-description)
  - [1.1 Architecture](#11-architecture)
  - [1.2 Noise Filter](#12-Noise-Filter)
  
  - [2. Installation](#2-installation)
  - [2.1 Docker](#21-docker)
    - [2.1.1 Docker Hub Image](#211-docker-hub-image)
  - [2.2 Local](#22-local)
- [3. Documentation](#3-documentation)
- [3. Easy Access](#3-easy-access)

## 1. Description

RIoTPot is an interoperable medium interaction honeypot, primarily focused on the emulation IoT and OT protocols, although, it is also capable of emulating other services.

This services are loaded in the honeypot in the form of plugins, making RIoTPot a modular, and very transportable honeypot. The services are loaded at runtime, meaning that the weight of the honeypot will vary on premisses, and the services loaded e.g. HTTP, will only be used when required. As consequence, we highly recommend building your own binary customized to your own needs. Refer to the following section, Installation, for more information.

### 1.1 Architecture

RIoTPot has a modular architecture that facilitates extensability of the honeypot. The honeypot further offers a hybrid-interaction capability where users can choosed the desired interaction levels for the protocols simulated. The image below shows the high/level architecture of RIoTPot. 

![alt text](architecture.jpg "Architecture")

The architecture contains 6 components. 

__RIoTPot core__
The core of the honeypot consists of the required modules for configuration, administration and orchestration of the container network.

__Configuration & Orchestration__
The configuration module provides RIoTPot with all the required parameters at startup. This includes the user preferences for specific protocols and profile simulation and the desired interaction level. The orchestration module is responsible for the network management from the core to the high-interaction protocol services simulated on containers. The received attack traffic is forwarded to the respective container that hosts the protocol on which the attack was targeted. Furthermore, the orchestra tor also facilitates the communication to the containers if they are hosted on a cloud-based environment.  

__Attack Capture and Noise Filter__
The  attack capture and  noise filter  module filters out the suspicious traffic received from Internet-wide scanners like Shodan and Censys. This helps the administrator to concentrate on attacks that are not from benign sources.

__Hybrid-Interaction (Low and High-Interaction modes)__
RIoTPot is implemented in Go language \cite{go} and facilitates the modular architecture and development through packages. The packages act as plug-ins that can be added to the honeypot to extend the protocols simulated. RIoTPot offers a hybrid-interaction model with a preference of low- or high-interaction. 
The low-interaction is achieved through independent packages, with each package simulating a specific protocol. The high-interaction model is realized with a containers with the protocols simulated as services installed. The containers act as high-interaction modules that offer a full implementation of the protocol. Additional protocol services can be added by integration of containers with desired protocol services. The hybrid-interaction model further allows the user to emulate selective protocols on low or high-interaction levels. For example, the user can choose to have SSH in low-interaction mode and MQTT in high-interaction mode thereby operating in a hybrid-interaction mode.

__Attack Database__
The attack database stores all the attack traffic received on the honeypot. The database is setup as an independent module to ensure data availability even if the honeypot crashes on potential large scale attacks. The database is accessible from the low-interaction and high-interaction modules for attack storage.

### 1.2 Noise Filter
The Noise filter module of RIoTPot filters the attacks from internet scanning engines to reduce alert fatigue.
With this feature, attacks are labelled as __*benign*__ when they originate from sources like Shodan. The 
list of scanning services filtered by RIoTPot is below:
 1. Shodan (https://www.shodan.io/)
 2. Censys (https://censys.io/)
 3. Project Sonar (https://www.rapid7.com/research/project-sonar/)
 4. LeakIX (https://leakix.net/)
 5. ShadowServer (https://www.shadowserver.org/)
 6. RWTH Aachen (http://researchscan.comsys.rwth-aachen.de/)
 7. Quadmetrics (https://www.quadmetrics.com/)
 8. BinaryEdge (https://www.binaryedge.io/})
 9. ipip.net (https://en.ipip.net/)
 10. Arbor Observatory (https://www.arbor-observatory.com/)
 11. CriminalIP (https://security.criminalip.com/)
 12. BitSight (https://www.bitsight.com/)
 13. InterneTT (http://www.internettl.org/)
 14. ONYPHE (https://www.onyphe.io/)
 15. Natlas (https://github.com/natlas/natlas)
 16. Net Systems Research (https://www.netsystemsresearch.com/)
 17. Sharashka (https://sharashka.io/data-feeds)
 18. Alpha Strike Labs (https://www.alphastrike.io)
 19. Stretchoid (http://stretchoid.com/)
 
 Note: the list will be updated on support for additional scanning sources. 

> **Summary:** To summarize, the design of RIoTPot facilitates modularity  through packages and containers as plugins. Furthermore, the modular architecture helps in achieving a hybrid-interaction model.

## 2. Installation

Although one can download the binaries and configuration files containing the set of default running emulators, this guide is mainly focused to those looking for a customized experience.

We thrive on the idea of making RIoTPot highly transportable, therefore, in this section one can find multiple methods of installation for diverse environments that fit a broad list of requirements and constrains.

We highly recommend running RiotPot in a virtualized self-contained network using `Docker`, for which we included configuration files that run the honeypot as a closed environment for testing and playing around (similar to a testbed environment).

> **NOTE:** The production image can be pulled from Docker Hub. If you choose this method you may directly jump to [2.1 Docker](#21-docker).
<!-- Not implemented yet, include section for the Docker Image pulling and building -->

RIoTPot is written in Golang, therefore, you will need to have [go](https://golang.org) installed first if you plan to make any changes, otherwise you can skip steps 1 and 2 if you rather not installing go. 

Regardless, you will need to copy RIoTPot to local:

```bash
# 1. Make the folder in where the repository will be stored.
$ mkdir -p $GOPATH/src/github.com

# 2. Navigate to the folder in where you store your repositories
$ mkdir -p $GOPATH/src/github.com

# 3. Clone the repository
$ git clone git@github.com:aau-network-security/riotpot.git

# 4. Navigate to the newly created folder with the repository
$ cd riotpot
```

### 2.1 Docker

We assume you have basic knowledge about the Docker ecosystem, otherwise please refer first to the Docker documentation [here](https://docs.docker.com).

At the deployments folder of RToTPot there is one **docker-compose** files:
```bash
$ cd ~/riotpot/deployments | ls -al
...
-rw-r--r--  docker-compose.yml
...
```

This file correspond to the respective software development environment *development*.

**Development.**
`docker-compose.yml` builds the project in a private virtual network in which there are three hosts: *riotpot*, *postgres*,and *tcpdump*. **Postgres** contains a postgres database, **tcpdump** contains a packet capturer, and **riotpot** the app itself. They can only communicate with the each other. Use this setup for **development** and **testing locally** by typing in your a terminal:

```bash
$ docker-compose -f docker-compose.yml up -d --build
```

Once you are done with the honeypot, you can put down the containers using the *down* command. 

```bash
$ docker-compose down -v
```

> **NOTE:** Using the *-v* tag will remove all the mounted volumes, i.e. the database used by riotpot to store information and the volumes mounted to store logs and binaries collected by the honeypot. Remember to make copies before using the *-v* tag, or skip it altogether.

#### 2.1.1 Docker Hub Image

Build the latest release of RiotPot directly from the image provided in the Docker Hub:

```bash
# Grab and run the latest release of the riotpot consumer image
# detached from the console with -d.
$ docker run -d riotpot-docker:latest
```

### 2.2 Local

To build your own binary from source, navigate to the folder where you have stored the repository and use the go CLI to generate it and store it in the *./bin/* folder:

```bash
# build the binary in the ./bin folder
$ go build -o riotpot cmd/riotpot/main.go
```

Additionally, you could also install the application in the system:
```bash
# installs riotpot at $GOPATH/bin
$ go install
```

Run the binary as any other application:
```bash
$ ./riotpot
```

## 3. Documentation

The documentation for RiotPot can be found in [go.pkg.dev](https://pkg.go.dev/), however, sometimes you might be in need to visualize the documentation locally, either because you are developing a part of it, of for any other reason.

The most common way of pre-visualizing documentation is by using `godoc`, however, this requires an initial setup of the go project. Find more information in the [godoc page](https://pkg.go.dev/golang.org/x/tools/cmd/godoc).

For simplicity, the riotpot `godoc` documentation can be run as a separated local container from the dockerfile `Dockerfile.documentation`. To use the container simply type:

```bash
$ make riotpot-doc
```
This will run a container tagged with `riotpot/v1` at `http://localhost:6060/`. The documentation of the package can be accessed directly from [http://localhost:6060/pkg/riotpot/](http://localhost:6060/pkg/riotpot/).

## 3. Easy Access

We previously described how to set up the whole project, both installation and documentation, but some of the processes become routinely and lengthy when on the process of developing new features and testing. For this, in the root folder of the repository we have included a `Makefile` containing the most utilized routines with aliases.

The following commands will be run using `make` plus the alias of the command. The `Makefile` contains more commands, but this are the most widely useful:

Command|Container Name|Description
:---|:---:|---:
riotpot-up|riotpot:development| Puts up RIoTPot in **development** mode.
riotpot-down|riotpot:development| Puts down RIoTPot.
riotpot-doc|riotpot/v1| Puts up a container with the local documentation.
riotpot-all|riotpot/v1, riotpot| Puts the documentation and RIoTPot **development** mode up.
riotpot-builder|| Builds the binary and the plugins.

**Example usage:**
```bash
# run a command given its alias from Makefile
$ make riotpot-doc
```
