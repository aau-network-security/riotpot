// This package implements an MQTT 3.1 honeypot
package main

import (
	"net"
	"sync"

	"github.com/riotpot/internal/globals"
	"github.com/riotpot/internal/logger"
	"github.com/riotpot/internal/services"
)

var Plugin string

const (
	name    = "MQTT"
	network = globals.TCP
	port    = 1883
)

func init() {
	Plugin = "Mqttd"
}

func Mqttd() services.Service {
	mx := services.NewPluginService(name, port, network)

	return &Mqtt{
		mx,
		sync.WaitGroup{},
	}
}

type Mqtt struct {
	services.Service
	wg sync.WaitGroup
}

func (m *Mqtt) Run() (err error) {
	// start a service in the `mqtt` port
	listener, err := net.Listen(m.GetNetwork().String(), m.GetAddress())
	logger.Log.Error().Err(err)

	// build a channel stack to receive connections to the service
	conn := make(chan net.Conn)

	// add a waiting group to serve the connections before continuing
	m.wg.Add(1)
	go m.serve(conn, listener)

	// handle the connections from the channel
	m.handlePool(conn)
	m.wg.Wait()

	return
}

// This function only serves to typical tcp connections, it currently does not handle
// websockets!!
func (m *Mqtt) serve(ch chan net.Conn, listener net.Listener) {
	defer m.wg.Done()

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

func (m *Mqtt) handlePool(ch chan net.Conn) {
	// open an infinite loop to handle the connections
	for {
		// while the `stop` channel remains empty, continue handling
		// new connections.
		select {
		case conn := <-ch:
			// use one goroutine per connection.
			go m.handleConn(conn)
		}
	}
}

func (m *Mqtt) handleConn(conn net.Conn) {
	// close the connection when the loop returns
	defer conn.Close()

	// Create a session for the connection
	// TODO include a list of topics as default that the
	// client can subscribe to.
	s := NewSession(conn)

	for {
		// read the connection packet
		packet := s.Read(conn)
		if packet == nil {
			// close the connection if the header is empty
			return
		}

		// respond to the message
		s.Answer(*packet, &conn)
	}

}
