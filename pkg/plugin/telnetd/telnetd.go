package main

import (
	"bufio"
	"fmt"
	"io/ioutil"
	"log"
	"net"

	"github.com/riotpot/internal/services"
	"github.com/riotpot/pkg/fake/shell"
	"github.com/riotpot/tools/errors"
)

var Plugin string

var (
	name     = "Telnetd"
	protocol = "tcp"
	port     = 23
	host     = "localhost"
)

func init() {
	Plugin = name
}

func Telnetd() services.Service {
	mx := services.NewPluginService(name, port, protocol, host)
	content, err := ioutil.ReadFile("banner.txt")
	if err != nil {
		log.Fatal(err)
	}

	return &Telnet{
		Service: mx,
		banner:  content,
	}
}

type Telnet struct {
	services.Service
	banner []byte
}

func (t *Telnet) Run() (err error) {

	// convert the port number to a string that we can use it in the server
	var port = fmt.Sprintf(":%d", t.GetPort())

	// start a service in the `telnet` port
	listener, err := net.Listen(t.GetProtocol(), port)
	errors.Raise(err)

	// build a channel stack to receive connections to the service
	conn := make(chan net.Conn)
	go t.serve(conn, listener)

	// handle the connections from the channel
	t.handlePool(conn)

	return
}

func (t *Telnet) serve(ch chan net.Conn, listener net.Listener) {
	// open an infinite loop to receive connections
	fmt.Printf("[%s] Started listenning for connections in port %d\n", t.GetName(), t.GetPort())
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
	for conn := range ch {
		go t.handleConn(conn)
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
	t.respond(string(t.banner), conn, br)

	pass := `Password: `
	t.respond(pass, conn, br)
}

// Offers a telnet shell-like experience in where
// the client will be prompt for input and the commands
// will be saved in the database.
func (t *Telnet) telnetShell(conn net.Conn, br *bufio.Reader) {
	// load a unix-like fake shell
	shell := shell.New("root", "ubuntu")
	shell.SetIo(conn)
	shell.Start()
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

	return
}
