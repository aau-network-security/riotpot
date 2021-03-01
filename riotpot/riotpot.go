// Main Application
package main

import (
	"log"
	"sync"
	"os"
	"testing"
	// we already have all of this locally
	//"github.com/sastry17/riotpot/internal/emulator/httpd"
	//"github.com/sastry17/riotpot/internal/emulator/sshd"
	//"github.com/sastry17/riotpot/internal/emulator/telnetd"
	//"github.com/sastry17/riotpot/external/trudy"
)


var wg sync.WaitGroup


func main() {

	if os.Getenv("ALL") {
		wg.Add(1)
		go telnet_serv()
		go http_serv()
		go ssh_serv()  //Starts SSH Server
		go start_proxy()
   		wg.Wait()
	}else{

	}
}

func serv_out(service string, out int){
	""" 
	Print to the console if the server successfully started
	@service: the name of the service
	@out: the integer result of the server started, 1 or 0 , where 0 is a falsy value.
	"""

	if out{
		fmt.Printf("%s Server Started", service)
	}else{
		fmt.Printf("[ERROR] Something went wrong during %s Server initialization", service)
	}
}


func ssh_serv(){
	"""Start SSH server"""
	var ssh_server int = go sshd.SSHServer()
	serv_out("SSH", ssh_server)
}


func telnet_serv(){
	"""Start Telnet server"""
	var telnet_server int = telnetd.TelnetServer() 
	serv_out("Telnet", telnet_server)
}


func http_serv(){
	"""Start HTTP server"""
	var http_server int = httpd.HttpServer()
	serv_out("HTTP", http_server)
}


func trydy_proxy_serv(){
	"""Start Trudy Proxy server"""
	var tryd_proxy int = trudy.Trudy()
	serv_out("Trudy Proxy", tryd_proxy)
}
