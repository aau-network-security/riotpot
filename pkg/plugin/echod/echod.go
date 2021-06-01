package main

import (
	"bufio"
	"fmt"
	"net"

	"github.com/riotpot/pkg/models"
	"github.com/riotpot/pkg/services"
	"github.com/riotpot/tools/errors"
)

var Name string

func init() {
	Name = "Echod"
}

func Echod() services.Service {
	mixin := services.MixinService{
		Name:     Name,
		Port:     7,
		Protocol: "tcp",
		Running:  make(chan bool, 1),
	}

	return &Echo{
		mixin,
	}
}

type Echo struct {
	// Anonymous fields from the mixin
	services.MixinService
}

func (e *Echo) Run() (err error) {
	// before running, migrate the model that we want to store
	e.Migrate(&models.Connection{})

	// convert the port number to a string that we can use it in the server
	var port = fmt.Sprintf(":%d", e.Port)

	// start a service in the `echo` port
	listener, err := net.Listen(e.Protocol, port)
	errors.Raise(err)

	// create the channel for stopping the service
	e.StopCh = make(chan int, 1)

	// build a channel stack to receive connections to the service
	conn := make(chan net.Conn)
	go e.serve(conn, listener)

	// update the status of the service
	e.Running <- true

	// handle the connections from the channel
	e.handlePool(conn)

	// Close the channel for stopping the service
	fmt.Print("[x] Service stopped...\n")
	close(e.StopCh)

	return
}

// Open the service and listen for connections
// inspired on https://gist.github.com/paulsmith/775764#file-echo-go
func (e *Echo) serve(ch chan net.Conn, listener net.Listener) {
	// open an infinite loop to receive connections
	fmt.Printf("[%s] Started listenning for connections in port %d\n", Name, e.Port)
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

// Handle the pool of connections to the service
func (e *Echo) handlePool(ch chan net.Conn) {
	// open an infinite loop to handle the connections
	for {
		// while the `stop` channel remains empty, continue handling
		// new connections.
		select {
		case <-e.StopCh:
			// stop the pool
			fmt.Printf("[x] Stopping %s service...\n", e.Name)
			// update the status of the service
			e.Running <- false
			return
		case conn := <-ch:
			// use one goroutine per connection.
			go e.handleConn(conn)
		}
	}
}

// Handle a connection made to the service
func (e *Echo) handleConn(conn net.Conn) {
	//opens a new small buffer
	br := bufio.NewReader(conn)

	for {
		// Read the message sent from the client.
		msg, err := br.ReadBytes('\n')
		if err != nil { // EOF, or worse
			break
		}

		// save the connection in the database
		e.save(conn, msg)
		// Respond with the same message
		conn.Write(msg)
	}
}

func (e *Echo) save(conn net.Conn, payload []byte) {

	connection := models.NewConnection()
	connection.LocalAddress = conn.LocalAddr().String()
	connection.RemoteAddress = conn.RemoteAddr().String()
	connection.Protocol = "TCP"
	connection.Service = Name
	connection.Incoming = true
	connection.Payload = string(payload)

	e.Store(connection)
}
