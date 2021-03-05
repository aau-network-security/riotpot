<div style="text-align:center"> 
  <img src="AAUgraphics/aau_logo.png" height="100">
</div> 

# Riotpot

<p align="left">
  <a href="https://github.com/aau-network-security/riotpot/actions"><img alt="GitHub Actions status" src="https://github.com/aau-network-security/riotpot/workflows/cyber/badge.svg"></a> 
</p>

# Table of Contents
- [Riotpot](#riotpot)
- [Table of Contents](#table-of-contents)
- [1. Description](#1-description)
- [2. Installation](#2-installation)


# 1. Description
Riotpot is an interoperable medium interaction honeypot, primarly focused on the emulation IoT and OT protocols, althoug, it is also capable of emulating other services.

This services are loaded in the honeypot in the form of plugins, making Riotpot a modular, and very transportable honeypot. The services are loaded at runtime, meaning that the weight of the honeypot will vary on premisses, and the services loaded e.g. HTTP, will only be used when required.

# 2. Installation

Docker container:
```bash
docker-compose up -d --build
```

Bare metal:
```bash
go build -o bin/
```

Specific file testing (e.g. plugins):
```bash
go run Path
```

<style> 
  body{ 
    background-image: url("AAUgraphics/aau_waves.png")
  } 
</style> 