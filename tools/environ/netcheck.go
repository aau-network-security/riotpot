/*
Package environ provides functions used to interact with the environment
*/
package environ

import (
	"fmt"
	"net"
)

/*
	Check if the port on the host machine is busy or not
	this is used for plugins to play on the host
*/
func CheckPortBusy(Protocol string, Port string) bool {
	_, err := net.Listen(Protocol, Port)

	if err != nil {
		fmt.Printf("Port %q is busy\n", Port)
		return false
	}

	return true
}

func CheckIPAddress(ip string) bool {
    if net.ParseIP(ip) == nil {
        fmt.Printf("IP Address: %s - Invalid\n", ip)
        return false
    } else {
        fmt.Printf("IP Address: %s - Valid\n", ip)
        return true
    }
}