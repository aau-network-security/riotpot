// Main Application
package main

import (
	"fmt"
	"sync"

	"riotpot/emulators/httpd"
	"riotpot/emulators/sshd"
	"riotpot/emulators/telnetd"
	"riotpot/utils/environ"
)

var wg sync.WaitGroup

// main starts all the submodules containing the emulators.
// It is the first function called when the application is run.
// It also acts as an orchestrator, which dictates the functioning of the application.
func main() {
	if environ.Getenv("SERVICES", "") == "" {
		wg.Add(1) // add a goroutine to the stack, blocking until all services are up.

		//----------------
		// Add here all the services that must be run by default if either there is no
		// environment file or variables which define this behaviour otherwise.
		//----------------
		go telnet_serv()
		go http_serv()
		go ssh_serv()
		//go trydy_proxy_serv()
		//----------------

		wg.Wait() // wait until the goroutines call Done or are over.
	} else {
		// TODO: add the logic to the loop.
	}
}

/*
Print to the console if the server successfully started
@service: the name of the service
@out: the integer result of the server started, 1 or 0 , where 0 is a falsy value.
*/
func serv_out(service string, out int) {

	if out > 0 {
		fmt.Printf("%s Server Started", service)
	} else {
		fmt.Printf("[ERROR] Something went wrong during %s Server initialization", service)
	}
}

/*Start SSH server*/
func ssh_serv() {
	sshd.SSHServer()
}

/* Start Telnet server */
func telnet_serv() {
	telnetd.TelnetServer()
}

/* Start HTTP server */
func http_serv() {
	httpd.HttpServer()
}

//func trydy_proxy_serv() {
/* Start Trudy Proxy server */
//	trudy.Trudy()
//}
