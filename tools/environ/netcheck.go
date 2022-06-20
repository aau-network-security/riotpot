/*
Package environ provides functions used to interact with the environment
*/
package environ

import (
	"log"
	"net"
)

/*
	Check if the port on the host machine is busy or not
	this is used for plugins to play on the host
*/
// TODO: change the hardcoded address of the host. At the very least get it from the configuration
func CheckPortBusy(protocol string, port string) bool {
	// Whether we can connect to some port
	conn, err := net.Listen(protocol, "localhost:"+port)
	isbusy := err != nil

	if isbusy {
		conn.Close()
	}

	return isbusy
}

// check if the IP address is valid
func CheckIPAddress(ip string) bool {
	isvalid := net.ParseIP(ip) != nil

	if !isvalid {
		log.Fatalf("Invalid IP %s", ip)
	}

	return isvalid
}

// check if IP address is reachable via ping command
func CheckIPConnection(IP string) {
	path := GetPath("ping")
	ExecuteCmd(path, IP, "-c", "2")
}
