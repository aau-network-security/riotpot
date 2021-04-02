package main

import (
	"bufio"
	"fmt"
	"net"
	"riotpot/services"
	"riotpot/utils/errors"
)

var Name string

func init() {
	Name = "Echod"
}

func Echod() services.Service {
	return &Echo{
		name: Name,
		// the natural port for `echo` is 7
		port: 7,
	}
}

type Echo struct {
	id   int
	name string
	port int
	stop chan int
}

func (e *Echo) Init(map[string]interface{}) {}

func (e *Echo) Run() error {
	var err error

	// convert the port number to a string that we can use in the server
	var port = fmt.Sprintf(":%b", e.port)

	// start a service in the `echo` port
	listener, err := net.Listen("tcp", port)
	errors.Raise(err)

	// build a channel stack to receive connections to the service
	conn := make(chan net.Conn)
	go e.serve(conn, listener)

	// handle the connections from the channel
	e.handlePool(conn)

	return err
}

// Open the service and listen for connections
// inspired on https://gist.github.com/paulsmith/775764#file-echo-go
func (e *Echo) serve(ch chan net.Conn, listener net.Listener) {
	// open an infinite loop to receive connections
	for {
		// Accept the client connection
		client, err := listener.Accept()
		errors.Raise(err)

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
		case conn := <-ch:
			// use one goroutine per connection.
			go e.handleConn(conn)
		case <-e.stop:
			fmt.Printf("[x] Stopping %s service...", e.name)
			return
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
		errors.Raise(err)

		// save the connection in the database
		e.save(conn.RemoteAddr(), msg)
		// Respond with the same message
		conn.Write(msg)
	}
}

func (e *Echo) save(address net.Addr, msg []byte) {
	// TODO: add logic here to store the connections made
	// to the server
	fmt.Printf("Bip Bop... \nSaving connection:\n> address: %v\n> msg: %v", address, msg)
}

func (e *Echo) Stop() error {
	var err error
	// send a stop signal to the channel
	e.stop <- 1
	return err
}

func (e *Echo) Restart() error {
	var err error
	return err
}

func (e *Echo) Status() error {
	var err error
	return err
}

func (e *Echo) Logger(ch chan<- error) (services.Logger, error) {
	var (
		logger services.Logger
		err    error
	)
	return logger, err
}
