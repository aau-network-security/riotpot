/*
Package environ provides functions used to interact with the environment
*/
package environ

import (
	"net"
	"log"
)

/*
	Check if the port on the host machine is busy or not
	this is used for plugins to play on the host
*/
func CheckPortBusy(protocol string, port string) bool {
	conn, err := net.Listen(protocol, "localhost:"+port)

	if err != nil {
		return false
	}

	conn.Close()
	return true
}

// check if the IP address is valid
func CheckIPAddress(ip string) bool {
    if net.ParseIP(ip) == nil {
        log.Fatalf("IP Address: %s - Invalid\n", ip)
        return false
    } else {
        return true
    }
}

// check if IP address is reachable via ping command
func CheckIPConnection(IP string) {
	path := GetPath("ping")
	ExecuteCmd(path, IP, "-c", "2")
}