// This package includes multiple validator implementations
package validators

import (
	"fmt"
	"net"
)

// Returns whether a port number is valid
func ValidatePortNumber(port int) (err error) {
	// Check if there is a port and is acceptable
	if !(port < 65536 && port > 0) {
		err = fmt.Errorf("invalid port %d", port)
		return
	}
	return
}

// Returns whether the port is available
// [7/18/2022] NOTE: this only considers TCP!!
func ValidatePortAvailable(port int) (err error) {
	// Check if the port is taken
	ln, err := net.Listen("tcp", fmt.Sprintf(":%d", port))
	if err != nil {
		return
	}
	defer ln.Close()
	return
}

// Wrapper to hecks whether the port is a valid number and available
// [7/18/2022] NOTE: this only considers TCP!!
func ValidatePort(port int) (p int, err error) {
	// Check if there is a port and is acceptable
	err = ValidatePortNumber(port)
	if err != nil {
		return
	}

	// Check if the port is available
	err = ValidatePortAvailable(port)
	if err != nil {
		return
	}

	p = port
	return
}
