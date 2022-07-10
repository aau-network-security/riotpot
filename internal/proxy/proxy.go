package proxy

import (
	"fmt"
	"net"
	"sync"

	"github.com/riotpot/pkg/services"
)

const (
	TCP = "tcp"
	UDP = "udp"
)

// Proxy interface.
type Proxy interface {
	// Start proxy function
	Start()

	// Stop the proxy
	Stop() error
	// Check if the proxy is running
	Alive() bool

	// Setter and Getter for the port
	SetPort(port int) (int, error)
	GetPort() int

	// Set the service in the proxy
	SetService(service services.Service)
}

// Abstraction of the proxy endpoint
// Contains private fields, do not use outside of this package
type AbstractProxy struct {
	Proxy

	// Port in where the proxy will listen
	port int
	// Protocol meant for this proxy
	protocol string

	// Create a channel to stop the proxy gracefully
	// This channel is also used to guess if the proxy is running
	stop chan struct{}

	// Pointer to the slice of middlewares for the proxies
	// All the proxies should apply and share the same middlewares
	// Perhaps this can be changed in the future given the need to apply middlewares per proxy
	middlewares *MiddlewareManagerItem

	// Service to proxy
	service services.Service

	// Waiting group for the server
	wg sync.WaitGroup
}

// Simple function to check if the proxy is running
func (pe *AbstractProxy) Alive() (alive bool) {
	// When the proxy is instantiated, the stop channel is nil;
	// therefore, the proxy is not running
	if pe.stop == nil {
		return
	}

	// [7/4/2022] NOTE: The logic of this block is difficult to read.
	// However, the select block will only give the default value when there is nothing
	// to read from the channel while the channel is still open.
	// When the channel is closed, the first case is not blocked, so we can not
	// read "anything else" from the channel
	select {
	// Return if the channel is closed
	case <-pe.stop:
	// Return if the channel is open
	default:
		alive = true
	}

	return
}

// Set the port based on some criteria
func (pe *AbstractProxy) SetPort(port int) (p int, err error) {
	p = port
	// Check if there is a port and is acceptable
	if !(port < 65536 && port > 0) {
		err = fmt.Errorf("invalid port %d", port)
		return
	}

	// Check if the port is taken
	ln, err := net.Listen(pe.protocol, fmt.Sprintf(":%d", port))
	if err != nil {
		return
	}
	defer ln.Close()

	pe.port = port
	return
}

// Returns the proxy port
func (pe *AbstractProxy) GetPort() int {
	return pe.port
}

func (pe *AbstractProxy) SetService(service services.Service) {
	pe.service = service
}

// Create a new instance of the proxy
func NewProxyEndpoint(port int, protocol string) (pe Proxy, err error) {
	// Get the proxy for UDP or TCP
	switch protocol {
	case TCP:
		pe, err = NewTCPProxy(port)
	case UDP:
		pe, err = NewUDPProxy(port)
	}

	return
}

// Interface for the proxy manager
type ProxyManager interface {
	// Create a new proxy and add it to the manager
	CreateProxy(port int) (*TCPProxy, error)
	// Delete a proxy from the list
	DeleteProxy(port int) error
	// Get a proxy by the port it uses
	GetProxy(port int) (*TCPProxy, error)
	// Set the service for a proxy
	SetService(port int, service services.Service) (pe *TCPProxy, err error)
}

// Simple implementation of the proxy manager
// This manager has access to the proxy endpoints registered. However, it does not observe newly
//
type ProxyManagerItem struct {
	ProxyManager

	// List of proxy endpoints registered in the manager
	proxies []Proxy

	// Instance of the middleware manager
	middlewares *MiddlewareManagerItem
}

func (pm *ProxyManagerItem) CreateProxy(protocol string, port int) (pe Proxy, err error) {
	// Check if there is another proxy with the same port
	if proxy, _ := pm.GetProxy(port); proxy != nil {
		err = fmt.Errorf("proxy already registered")
		return
	}

	// Create the proxy
	pe, err = NewProxyEndpoint(port, protocol)

	// Append the proxy to the list
	pm.proxies = append(pm.proxies, pe)
	return
}

// Delete a proxy from the registered list
// The proxy is stopped before being removed
func (pm *ProxyManagerItem) DeleteProxy(port int) (err error) {
	// Iterate the registered proxies for the proxy using the given port, and stop and remove it from the slice
	for ind, proxy := range pm.proxies {
		if proxy.GetPort() == port {
			// Stop the proxy, just in case
			proxy.Stop()
			// Remove it from the slice by replacing it with the last item from the slice, and reducing the slice
			// by 1 element
			lastInd := len(pm.proxies) - 1

			pm.proxies[ind] = pm.proxies[lastInd]
			pm.proxies = pm.proxies[:lastInd]
			return
		}
	}

	// If the proxy was not foun, send an error
	err = fmt.Errorf("proxy not found")
	return
}

// Returns a proxy by the port number
func (pm *ProxyManagerItem) GetProxy(port int) (pe Proxy, err error) {
	// Iterate the proxies registered, and if the proxy using the given port is found, return it
	for _, proxy := range pm.proxies {
		if proxy.GetPort() == port {
			pe = proxy
			return
		}
	}

	// If the proxy was not foun, send an error
	err = fmt.Errorf("proxy not found")
	return
}

// Set the service for some proxy
func (pm *ProxyManagerItem) SetService(port int, service services.Service) (pe Proxy, err error) {
	// Get the proxy from the list
	pe, err = pm.GetProxy(port)
	if err != nil {
		return
	}

	// If the proxy was found, set the service
	pe.SetService(service)

	return
}

// Constructor for the proxy manager
func NewProxyManager() *ProxyManagerItem {
	return &ProxyManagerItem{
		middlewares: Middlewares,
		proxies:     make([]Proxy, 0),
	}
}
