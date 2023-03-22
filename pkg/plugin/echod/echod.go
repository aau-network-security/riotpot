package main

import (
	"bufio"
	"net"

	"github.com/riotpot/internal/globals"
	"github.com/riotpot/internal/logger"
	"github.com/riotpot/internal/services"
)

var Plugin string

const (
	name    = "Echo"
	port    = 7
	network = globals.TCP
)

func init() {
	Plugin = "Echod"
}

func Echod() services.Service {
	mx := services.NewPluginService(name, port, network)

	return &Echo{
		mx,
	}
}

type Echo struct {
	// Anonymous fields from the mixin
	services.Service
}

func (e *Echo) Run() (err error) {
	// start a service in the `echo` port
	listener, err := net.Listen(e.GetNetwork().String(), e.GetAddress())
	logger.Log.Error().Err(err)

	// build a channel stack to receive connections to the service
	conn := make(chan net.Conn)
	go e.serve(conn, listener)

	// handle the connections from the channel
	e.handlePool(conn)

	return
}

// Open the service and listen for connections
// inspired on https://gist.github.com/paulsmith/775764#file-echo-go
func (e *Echo) serve(ch chan net.Conn, listener net.Listener) {
	// open an infinite loop to receive connections
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

		// Respond with the same message
		conn.Write(msg)
	}
}
