package main

import (
	"fmt"
	"github.com/sastry17/riotpot/internal/emulator/httpd"
	"github.com/sastry17/riotpot/internal/emulator/sshd"
	"github.com/sastry17/riotpot/internal/emulator/telnetd"
	"sync"
	"github.com/sastry17/riotpot/external/trudy"
	"os"
)
var wg sync.WaitGroup
func main() {

	args := os.Args
	if args[1] == "--all" {

		wg.Add(1)
		go telnet_serv()
		go http_serv()
		fmt.Println("Telnet and HTTP Server Started")
		go sshd.SSHServer()  //Starts SSH Server
		fmt.Println("SSH Server Started")
		go start_proxy()
   		wg.Wait()

	}

}

func telnet_serv(){
	telnetd.TelnetServer() //Starts Telnet emulation
}

func http_serv() {
	httpd.HttpServer() //starts HTTP Server
}

func start_proxy() {
	trudy.Trudy() //start trudy proxy
}