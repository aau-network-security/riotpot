package proxy

import (
	"fmt"
	"io"
	"net"
	"sync"
	"time"

	"github.com/riotpot/internal/globals"
	lr "github.com/riotpot/internal/logger"
)

// Implementation of a TCP proxy

type TCPProxy struct {
	*AbstractProxy
	listener net.Listener
}

// Start listening for connections
func (tcpProxy *TCPProxy) Start() (err error) {
	// Check if the service is set, otherwise return with an error
	if tcpProxy.GetService() == nil {
		err = fmt.Errorf("service not set")
		return
	}

	// Get the listener or create a new one
	listener, err := tcpProxy.GetListener()

	if err != nil {
		return
	}

	// Create a channel to stop the proxy
	tcpProxy.stop = make(chan struct{})

	// Add a waiting task
	tcpProxy.wg.Add(1)

	go func() {
		defer tcpProxy.wg.Done()

		for {
			// Accept the next connection
			// This goes first as it is the method we have to check if the proxy is running
			// There is no need to continue if it is not
			client, err := listener.Accept()
			if err != nil {
				return
			}
			defer client.Close()

			// Apply the middlewares to the connection before dialing the server
			_, err = tcpProxy.middlewares.Apply(client)
			if err != nil {
				return
			}

			// Get a connection to the server for each new connection with the client
			server, servErr := net.DialTimeout(globals.TCP.String(), tcpProxy.service.GetAddress(), 1*time.Second)

			// If there was an error, close the connection to the server and return
			if servErr != nil {
				return
			}
			defer server.Close()

			// Add a waiting task
			tcpProxy.wg.Add(1)

			go func() {
				// Handle the connection between the client and the server
				// NOTE: The handlers will defer the connections
				tcpProxy.handle(client, server)

				// Finish the task
				tcpProxy.wg.Done()
			}()
		}
	}()

	return
}

func (tcpProxy *TCPProxy) GetListener() (listener net.Listener, err error) {
	listener = tcpProxy.listener

	// Get the listener only
	if listener == nil || tcpProxy.GetStatus() != globals.RunningStatus {
		listener, err = tcpProxy.NewListener()
		if err != nil {
			return
		}
		tcpProxy.listener = listener
	}

	return
}

func (tcpProxy *TCPProxy) NewListener() (listener net.Listener, err error) {
	listener, err = net.Listen(tcpProxy.GetNetwork().String(), fmt.Sprintf(":%d", tcpProxy.GetPort()))
	tcpProxy.AbstractProxy.listener = listener
	return
}

// TCP synchronous tunnel that forwards requests from source to destination and back
func (tcpProxy *TCPProxy) handle(from net.Conn, to net.Conn) {
	// Create the waiting group for the connections so they can answer the each other
	var wg sync.WaitGroup
	wg.Add(2)

	handler := func(source net.Conn, dest net.Conn) {
		defer wg.Done()

		// Write the content from the source to the destination
		_, err := io.Copy(dest, source)
		if err != nil {
			lr.Log.Warn().Err(err).Msg("Could not copy from source to destination")
		}

		// Close the connection to the source
		if err := source.Close(); err != nil {
			lr.Log.Warn().Err(err)
		}

		// Attempt to close the writter. This may not always work
		// Another solution is to just call `Close()` on the writter
		if d, ok := dest.(*net.TCPConn); ok {
			if err := d.CloseWrite(); err != nil {
				lr.Log.Warn().Err(err)
			}

		}
	}

	// Start the workers
	// TODO: [7/3/2022] Check somewhere if the connection is still alive from the source and destination
	// Otherwise there is no need to wait
	go handler(from, to)
	go handler(to, from)

	// Wait until the forwarding is done
	wg.Wait()
}

func NewTCPProxy(port int) (proxy *TCPProxy, err error) {
	// Create a new proxy
	proxy = &TCPProxy{
		AbstractProxy: NewAbstractProxy(port, globals.TCP),
	}

	// Set the port
	proxy.SafeSetPort(port)
	return
}
