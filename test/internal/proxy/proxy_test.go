package proxy

import (
	"fmt"
	"io/ioutil"
	"net"
	"testing"

	"github.com/riotpot/internal/globals"
	"github.com/riotpot/internal/proxy"
	"github.com/riotpot/internal/services"
	"github.com/stretchr/testify/assert"
)

const (
	proxyPort  = 8080
	serverPort = 8081
	network    = globals.TCP
)

// Test to create a service, the ports at the end should be the same
func TestCreateProxy(t *testing.T) {
	// Instantiate the proxy
	pr, err := proxy.NewProxyEndpoint(proxyPort, network)

	if err != nil {
		t.Error(err)
	}

	assert.Equal(t, pr.GetPort(), proxyPort, "The ports should be the same")
}

// Test to start the proxy
// It creates a local client and server. The message should be transferred from
// the client to the server trhough the proxy
func TestStartProxy(t *testing.T) {
	assert := assert.New(t)

	// use a channel to hand off the error
	errs := make(chan error, 1)

	// Message to send
	message := "Hi there!"
	ret := ""

	// Instantiate the proxy
	pr, err := proxy.NewProxyEndpoint(proxyPort, network)
	if err != nil {
		t.Fatal(err)
	}

	// Create a new abstract service
	service := services.NewService("http", serverPort, network, "", globals.Low)

	// Set the service
	pr.SetService(service)

	// Start a listener for that port and accept anything (server)
	go func() {
		l, err := net.Listen(network.String(), service.GetAddress())
		if err != nil {
			errs <- err
		}
		defer l.Close()

		// Accept the connection
		conn, err := l.Accept()
		if err != nil {
			errs <- err
		}
		defer conn.Close()

		buf, err := ioutil.ReadAll(conn)
		if err != nil {
			errs <- err
		}

		ret = string(buf[:])

		// Check that the first message we got is the same as the expected
		assert.Equal(ret, message, "The messages must be equal")

		// Close the connection and return
		l.Close()
	}()

	// Start the proxy
	pr.Start()

	// Start a new client connections that sends the message (client)
	go func() {
		// Connect to the proxy
		conn, err := net.Dial(network.String(), fmt.Sprintf(":%d", proxyPort))
		if err != nil {
			errs <- err
		}
		defer conn.Close()

		// Send the message through the connection
		if _, err := fmt.Fprint(conn, message); err != nil {
			errs <- err
		}
	}()

	// wait for it
	err = <-errs
	if err != nil {
		t.Fatal(err)
	}

	assert.Equal(ret, message, "The messages must be equal")
}

// Test to stop the proxy
func TestStopProxy(t *testing.T) {
	assert := assert.New(t)
	// use a channel to hand off the error
	errs := make(chan error, 1)

	// Instantiate the proxy
	pr, err := proxy.NewProxyEndpoint(proxyPort, network)
	if err != nil {
		t.Fatal(err)
	}

	// Create a new abstract service
	service := services.NewService("http", serverPort, network, "", globals.Low)

	// Set the service
	pr.SetService(service)

	// Start a listener for that port and accept anything (server)
	go func() {
		l, err := net.Listen(network.String(), service.GetAddress())
		if err != nil {
			errs <- err
		}
		defer l.Close()

		// Accept the connection
		conn, err := l.Accept()
		if err != nil {
			errs <- err
		}
		defer conn.Close()

		// Close the connection and return
		l.Close()
	}()

	// Start the proxy
	pr.Start()

	// Give the proxy some time to start
	alive := pr.GetStatus()
	assert.Equal(alive, true, "The proxy is running")

	// Stop the service
	pr.Stop()
	alive = pr.GetStatus()
	assert.Equal(alive, false, "The proxy is stop")
}

func TestCreateNewProxy(t *testing.T) {
	// Create a new proxy manager
	proxyManager := proxy.NewProxyManager()
	// Add a proxy
	_, err := proxyManager.CreateProxy(globals.TCP, proxyPort)

	// There would be an error if the proxy was already registered or the port is unavailable
	if err != nil {
		t.Fatal(err)
	}
}

func TestDeleteProxy(t *testing.T) {
	// Create a new proxy manager
	proxyManager := proxy.NewProxyManager()
	// Add a proxy
	pe, err := proxyManager.CreateProxy(globals.TCP, proxyPort)
	if err != nil {
		t.Fatal(err)
	}

	// Delete the proxy
	err = proxyManager.DeleteProxy(pe.GetID())
	// There may be an error if the proxy was not found
	if err != nil {
		t.Fatal(err)
	}
}
