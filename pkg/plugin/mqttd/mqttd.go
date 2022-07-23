// This package implements an MQTT 3.1 honeypot
package main

import (
	"fmt"
	"net"
	"sync"

	"github.com/riotpot/pkg/profiles/ports"
	"github.com/riotpot/pkg/services"
	"github.com/riotpot/tools/errors"
)

var Name string

func init() {
	Name = "Mqttd"
}

func Mqttd() services.Service {
	mx := services.NewPluginService(Name, ports.GetPort(Name), "tcp")

	return &Mqtt{
		mx,
		sync.WaitGroup{},
	}
}

type Mqtt struct {
	*services.PluginService
	wg sync.WaitGroup
}

func (m *Mqtt) Run() (err error) {

	// convert the port number to a string that we can use it in the server
	var port = fmt.Sprintf(":%d", m.GetPort())

	// start a service in the `mqtt` port
	listener, err := net.Listen(m.GetProtocol(), port)
	errors.Raise(err)

	// create the channel for stopping the service
	m.StopCh = make(chan int, 1)

	// build a channel stack to receive connections to the service
	conn := make(chan net.Conn)

	// add a waiting group to serve the connections before continuing
	m.wg.Add(1)
	go m.serve(conn, listener)

	// update the status of the service
	m.Running <- true

	// handle the connections from the channel
	m.handlePool(conn)

	// Close the channel for stopping the service
	fmt.Print("[x] Service stopped...\n")
	close(m.StopCh)

	return
}

// This function only serves to typical tcp connections, it currently does not handle
// websockets!!
func (m *Mqtt) serve(ch chan net.Conn, listener net.Listener) {
	defer m.wg.Done()

	// open an infinite loop to receive connections
	fmt.Printf("[%s] Started listenning for connections in port %d\n", Name, m.GetPort())
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
		case <-m.StopCh:
			// stop the pool
			fmt.Printf("[x] Stopping %s service...\n", m.GetName())
			// update the status of the service
			m.Running <- false
			return
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
