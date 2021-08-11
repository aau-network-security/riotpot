package ports

var ports = map[string]int {"coapd": 5683,
							"sshd": 22,
							"httpd": 8080,
							"echod": 7,
							"telnetd": 23,
							"mqttd": 1883,
							"modbusd": 502}

func GetPort(service string) (int) {
	return ports[service]
}
