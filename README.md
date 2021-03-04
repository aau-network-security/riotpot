# riotpot

<p align="left">
  <a href="https://github.com/aau-network-security/riotpot/actions"><img alt="GitHub Actions status" src="https://github.com/aau-network-security/riotpot/workflows/cyber/badge.svg"></a> 
</p>

Resilient IoT Hoeypot

# Emulators

Riotpot offers a diverse set of protocol emulators:
- HTTP
- SSH
- MQTT
- TELNET
- ECHO

# Usage

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