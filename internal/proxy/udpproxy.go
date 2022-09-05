package proxy

import (
	"fmt"
	"net"
	"sync"

	"github.com/riotpot/internal/globals"
	lr "github.com/riotpot/internal/logger"
)

type UDPProxy struct {
	*AbstractProxy
	listener *net.UDPConn
}

func (udpProxy *UDPProxy) Start() (err error) {
	// Check if the service is set
	if udpProxy.GetService() == nil {
		err = fmt.Errorf("service not set")
		return
	}

	// Get the listener or create a new one
	client, err := udpProxy.GetListener()
	if err != nil {
		return
	}
	defer client.Close()

	// Create a channel to stop the proxy
	udpProxy.stop = make(chan struct{})

	// Add a waiting task
	udpProxy.wg.Add(1)

	srvAddr := net.UDPAddr{
		Port: udpProxy.service.GetPort(),
	}

	for {
		// Get a connection to the server for each new connection with the client
		server, servErr := net.DialUDP(globals.UDP.String(), nil, &srvAddr)
		// If there was an error, close the connection to the server and return
		if servErr != nil {
			return
		}
		defer server.Close()

		go func() {
			// TODO: Handle the middlewares! they only accept TCP connections
			// Apply the middlewares to the connection
			//udpProxy.middlewares.Apply(listener)

			// Handle the connection between the client and the server
			// NOTE: The handlers will defer the connections
			udpProxy.handle(client, server)

			// Finish the task
			udpProxy.wg.Done()
		}()
	}
}

// Function to stop the proxy from runing
func (udpProxy *UDPProxy) Stop() (err error) {
	// Stop the proxy if it is still alive
	if udpProxy.GetStatus() != globals.StoppedStatus {
		close(udpProxy.stop)
		udpProxy.listener.Close()
		// Wait for all the connections and the server to stop
		udpProxy.wg.Wait()
		return
	}

	err = fmt.Errorf("proxy not running")
	return
}

// Get or create a new listener
func (udpProxy *UDPProxy) GetListener() (listener *net.UDPConn, err error) {
	listener = udpProxy.listener

	// Check if there is a listener
	if listener == nil || udpProxy.GetStatus() != globals.RunningStatus {
		// Get the address of the UDP server
		addr := net.UDPAddr{
			Port: udpProxy.service.GetPort(),
		}

		listener, err = net.ListenUDP(globals.UDP.String(), &addr)
		if err != nil {
			return
		}
		udpProxy.listener = listener
		udpProxy.AbstractProxy.listener = listener
	}

	return
}

// TODO: Test this function
// UDP asynchronous tunnel
func (udpProxy *UDPProxy) handle(client *net.UDPConn, server *net.UDPConn) {
	var buf [2 << 10]byte
	var wg sync.WaitGroup
	wg.Add(2)

	// Function to copy messages from one pipe to the other
	var handle = func(from *net.UDPConn, to *net.UDPConn) {
		n, addr, err := from.ReadFrom(buf[0:])
		if err != nil {
			lr.Log.Warn().Err(err)
		}

		_, err = to.WriteTo(buf[:n], addr)
		if err != nil {
			lr.Log.Warn().Err(err)
		}
	}

	defer client.Close()
	defer server.Close()

	go handle(client, server)
	go handle(server, client)

	// Wait until the forwarding is done
	wg.Wait()
}

func NewUDPProxy(port int) (proxy *UDPProxy, err error) {
	// Create a new proxy
	proxy = &UDPProxy{
		AbstractProxy: NewAbstractProxy(port, globals.UDP),
	}

	// Set the port
	proxy.SetPort(port)
	return
}
