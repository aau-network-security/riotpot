package proxy

import (
	"fmt"
	"io"
	"log"
	"net"
	"sync"
	"time"
)

// Implementation of a TCP proxy

type TCPProxy struct {
	*AbstractProxy
	listener net.Listener
}

// Start listening for connections
func (tcpProxy *TCPProxy) Start() {
	// Get the listener or create a new one
	listener := tcpProxy.GetListener()
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
				// If the channel was closed, the proxy should stop
				if !tcpProxy.Alive() {
					return
				}
				fmt.Println(err)
			}
			defer client.Close()

			// Get a connection to the server for each new connection with the client
			server, servErr := net.DialTimeout(TCP, tcpProxy.service.GetAddress(), 1*time.Second)

			// If there was an error, close the connection to the server and return
			if servErr != nil {
				server.Close()
				return
			}
			defer server.Close()

			// Add a waiting task
			tcpProxy.wg.Add(1)

			go func() {
				// Apply the middlewares to the connection
				tcpProxy.middlewares.Apply(client)

				// Handle the connection between the client and the server
				// NOTE: The handlers will defer the connections
				tcpProxy.handle(client, server)

				// Finish the task
				tcpProxy.wg.Done()
			}()
		}
	}()
}

// Function to stop the proxy from runing
func (tcpProxy *TCPProxy) Stop() (err error) {
	// Stop the proxy if it is still alive
	if tcpProxy.Alive() {
		close(tcpProxy.stop)
		tcpProxy.listener.Close()
		// Wait for all the connections and the server to stop
		tcpProxy.wg.Wait()
		return
	}

	err = fmt.Errorf("proxy not running")
	return
}

// Get or create a new listener
func (tcpProxy *TCPProxy) GetListener() net.Listener {
	if tcpProxy.listener == nil || !tcpProxy.Alive() {
		listener, err := net.Listen(tcpProxy.protocol, fmt.Sprintf(":%d", tcpProxy.GetPort()))
		if err != nil {
			log.Fatal(err)
		}
		tcpProxy.listener = listener
	}
	return tcpProxy.listener
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
			log.Print(err)
		}

		// Close the connection to the source
		if err := source.Close(); err != nil {
			log.Print(err)
		}

		// Attempt to close the writter. This may not always work
		// Another solution is to just call `Close()` on the writter
		if d, ok := dest.(*net.TCPConn); ok {
			if err := d.CloseWrite(); err != nil {
				log.Print(err)
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
		AbstractProxy: &AbstractProxy{
			middlewares: Middlewares,
			protocol:    TCP,
		},
	}

	// Set the port
	_, err = proxy.SetPort(port)
	return
}
