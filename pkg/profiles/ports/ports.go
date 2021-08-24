package ports

import (
	"strconv"
	"log"
)

var ProtocolDetails = map[string]map[string]string { "coapd": { "protocol": "udp", "port": "5683"},
													"sshd":  { "protocol": "tcp", "port": "22"},
													"httpd":  { "protocol": "tcp", "port": "8080"},
													"echod":  { "protocol": "tcp", "port": "7"},
													"telnetd":  { "protocol": "tcp", "port": "23"},
													"mqttd":  { "protocol": "tcp", "port": "1883"},
													"modbusd":  { "protocol": "tcp", "port": "502"},
													}

func GetPort(service string) (int) {
	// convert port from string to int
	val, err := strconv.Atoi(ProtocolDetails[service]["port"])

	if err != nil {
		log.Fatalf("Error occoured %q", err)
	}

	return val
}

func GetProtocol(service string) (string) {

	return ProtocolDetails[service]["protocol"]
}