package proxy

import (
	"fmt"
	"sync"

	"github.com/google/uuid"
	"github.com/riotpot/internal/validators"
	"github.com/riotpot/pkg/services"
)

const (
	// Protocols
	TCP = "tcp"
	UDP = "udp"
)

const (
	// Values for the proxy status
	// This variables are set to make development more readable and extensible
	// If in the future other status are allowed, include them here
	ALIVE = 1
	DEAD  = 0
)

// Proxy interface.
type Proxy interface {
	// Start and stop
	Start() error
	Stop() error

	// Getters
	GetID() string
	GetPort() int
	GetProtocol() string
	GetStatus() int
	GetService() services.Service

	// Setters
	SetPort(port int) int
	SetService(service services.Service) services.Service
}

// Abstraction of the proxy endpoint
// Contains private fields, do not use outside of this package
type AbstractProxy struct {
	Proxy

	// ID of the proxy
	id uuid.UUID

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

	// Generic listener
	listener interface{ Close() }
}

// Function to stop the proxy from runing
func (pe *AbstractProxy) Stop() (err error) {
	// Stop the proxy if it is still alive
	if pe.GetStatus() != DEAD {
		close(pe.stop)
		pe.listener.Close()
		// Wait for all the connections and the server to stop
		pe.wg.Wait()
		return
	}

	err = fmt.Errorf("proxy not running")
	return
}

// Simple function to check if the proxy is running
func (pe *AbstractProxy) GetStatus() (alive int) {
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
		alive = ALIVE
	}

	return
}

func (pe *AbstractProxy) GetID() string {
	return pe.id.String()
}

// Set the port
// NOTE: use the ValidatePort before assigning
func (pe *AbstractProxy) SafeSetPort(port int) (p int, err error) {
	p, err = validators.ValidatePort(port)
	if err != nil {
		return
	}

	pe.port = p
	return
}

// Set the port
// NOTE: use the ValidatePort before assigning
func (pe *AbstractProxy) SetPort(port int) int {

	pe.port = port
	return pe.port
}

// Returns the proxy port
func (pe *AbstractProxy) GetPort() int {
	return pe.port
}

// Set the service based on the list of registered services
func (pe *AbstractProxy) SetService(service services.Service) services.Service {
	pe.service = service
	return pe.service
}

// Returns the service
func (pe *AbstractProxy) GetService() services.Service {
	return pe.service
}

// Returns the service
func (pe *AbstractProxy) GetProtocol() string {
	return pe.protocol
}

func NewAbstractProxy(port int, protocol string) (ab *AbstractProxy) {
	ab = &AbstractProxy{
		id:          uuid.New(),
		port:        port,
		protocol:    protocol,
		middlewares: Middlewares,
	}
	return
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
