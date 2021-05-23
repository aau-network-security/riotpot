package main

import (
	"bufio"
	"fmt"
	"net"

	"github.com/riotpot/pkg/fake/shell"
	"github.com/riotpot/pkg/models"
	"github.com/riotpot/pkg/services"
	"github.com/riotpot/tools/errors"
)

var Name string

func init() {
	Name = "Telnetd"
}

func Telnetd() services.Service {
	mixin := services.MixinService{
		Name:     Name,
		Port:     23,
		Running:  make(chan bool, 1),
		Protocol: "tcp",
	}

	return &Telnet{
		mixin,
	}
}

type Telnet struct {
	services.MixinService
}

func (t *Telnet) Run() (err error) {
	// before running, migrate the model that we want to store
	t.Migrate(&models.Connection{})

	// convert the port number to a string that we can use it in the server
	var port = fmt.Sprintf(":%d", t.Port)

	// start a service in the `telnet` port
	listener, err := net.Listen(t.Protocol, port)
	errors.Raise(err)

	// create the channel for stopping the service
	t.StopCh = make(chan int, 1)

	// build a channel stack to receive connections to the service
	conn := make(chan net.Conn)
	go t.serve(conn, listener)

	// update the status of the service
	t.Running <- true

	// handle the connections from the channel
	t.handlePool(conn)

	// Close the channel for stopping the service
	fmt.Print("[x] Service stopped...\n")
	close(t.StopCh)

	return
}

func (t *Telnet) serve(ch chan net.Conn, listener net.Listener) {
	// open an infinite loop to receive connections
	fmt.Printf("[%s] Started listenning for connections in port %d\n", Name, t.Port)
	for {
		// Accept the client connection
		client, err := listener.Accept()
		if err != nil {
			return
		}
		defer client.Close()

		// push the client connection to the channel
		ch <- client
	}
}

func (t *Telnet) handlePool(ch chan net.Conn) {
	// open an infinite loop to handle the connections
	for {
		// while the `stop` channel remains empty, continue handling
		// new connections.
		select {
		case <-t.StopCh:
			// stop the pool
			fmt.Printf("[x] Stopping %s service...\n", t.Name)
			// update the status of the service
			t.Running <- false
			return
		case conn := <-ch:
			// use one goroutine per connection.
			go t.handleConn(conn)
		}
	}
}
func (t *Telnet) handleConn(conn net.Conn) {
	//opens a new small buffer
	br := bufio.NewReader(conn)

	// Send the authentication messages
	t.sendAuth(conn, br)
	// encarcelate the client in the telnet shell loop
	t.telnetShell(conn, br)
}

// This method shows the welcome message to the telnet
// service, and prompts for authentication.
func (t *Telnet) sendAuth(conn net.Conn, br *bufio.Reader) {
	// Welcome message
	hello := `
This device is for authorized personnel only.
If you have not been provided with permission to
access this device - disconnect at once.

*** Login Required.  Unauthorized use is prohibited ***
*** Ensure that you update the system configuration ***
*** documentation after making system changes.      ***

User Access Verification: `
	t.respond(hello, conn, br)

	pass := `Password: `
	t.respond(pass, conn, br)
}

// Offers a telnet shell-like experience in where
// the client will be prompt for input and the commands
// will be saved in the database.
func (t *Telnet) telnetShell(conn net.Conn, br *bufio.Reader) {
	// load a unix-like fake shell
	shell := shell.New()
	shell.SetIo(conn)
	shell.Start()

	go func() {
		for rq := range shell.RspChan {
			t.save(conn, rq)
		}
	}()
}

// Method to send a message to the client, receive a response and save it
// into the database.
func (t *Telnet) respond(
	msg string,
	conn net.Conn,
	br *bufio.Reader,
) (response []byte, err error) {
	// send the message and wait for the client to respond
	conn.Write([]byte(msg))

	// read the response
	response, err = br.ReadBytes('\n')
	if err != nil { // EOF, or worse
		return
	}

	// save the response in the database
	t.save(conn, response)
	return
}

func (t *Telnet) save(conn net.Conn, payload []byte) {

	connection := models.NewConnection()
	connection.LocalAddress = conn.LocalAddr().String()
	connection.RemoteAddress = conn.RemoteAddr().String()
	connection.Protocol = "TCP"
	connection.Service = Name
	connection.Incoming = true
	connection.Payload = string(payload)

	t.Store(connection)
}
