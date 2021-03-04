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

func main() {

	if environ.Getenv("SERVICES", "") == "" {
		wg.Add(1)
		go telnet_serv()
		go http_serv()
		go ssh_serv() //Starts SSH Server
		//go trydy_proxy_serv()
		wg.Wait()
	} else {

	}
}

func serv_out(service string, out int) {
	/*
		Print to the console if the server successfully started
		@service: the name of the service
		@out: the integer result of the server started, 1 or 0 , where 0 is a falsy value.
	*/

	if out > 0 {
		fmt.Printf("%s Server Started", service)
	} else {
		fmt.Printf("[ERROR] Something went wrong during %s Server initialization", service)
	}
}

func ssh_serv() {
	/*Start SSH server*/
	sshd.SSHServer()
}

func telnet_serv() {
	/* Start Telnet server */
	telnetd.TelnetServer()
}

func http_serv() {
	/* Start HTTP server */
	httpd.HttpServer()
}

//func trydy_proxy_serv() {
/* Start Trudy Proxy server */
//	trudy.Trudy()
//}
