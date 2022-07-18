// This package implements a proxy manager
// A proxy manager

package proxy

import (
	"fmt"

	"github.com/riotpot/pkg/services"
)

var (
	// Instantiate the proxy manager to allow other applications work with the proxies
	Proxies = NewProxyManager()
)

// Interface for the proxy manager
type ProxyManager interface {
	// Get all the proxies registered
	GetProxies() []Proxy
	// Create a new proxy and add it to the manager
	CreateProxy(protocol string, port int) (Proxy, error)

	// Methods for proxies the using ID field
	GetProxy(id string) (Proxy, error)
	DeleteProxy(id string) error

	// Wrapper method to find a proxy using the port and protocol
	GetProxyFromParams(protocol string, port int)

	// Set the service for a proxy
	SetService(port int, service services.Service) (pe Proxy, err error)
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

// Create a new proxy and add it to the manager
func (pm *ProxyManagerItem) CreateProxy(protocol string, port int) (pe Proxy, err error) {
	// Check if there is another proxy with the same port
	if proxy, _ := pm.GetProxyFromParams(protocol, port); proxy != nil {
		err = fmt.Errorf("proxy already registered")
		return
	}

	// Create the proxy
	pe, err = NewProxyEndpoint(port, protocol)

	// Append the proxy to the list
	pm.proxies = append(pm.proxies, pe)
	return
}

// Delete a proxy from teh registered list using the ID
func (pm *ProxyManagerItem) DeleteProxy(id string) (err error) {
	// Get all the proxies registered
	proxies := pm.GetProxies()

	for ind, proxy := range proxies {
		if proxy.GetID() == id {
			// Stop the proxy, just in case
			proxy.Stop()
			// Remove it from the slice by replacing it with the last item from the slice,
			// and reducing the slice by 1 element
			lastInd := len(proxies) - 1

			proxies[ind] = proxies[lastInd]
			pm.proxies = proxies[:lastInd]

			return
		}
	}
	return
}

func (pm *ProxyManagerItem) GetProxies() []Proxy {
	return pm.proxies
}

// Returns a proxy by the port number
func (pm *ProxyManagerItem) GetProxyFromParams(protocol string, port int) (pe Proxy, err error) {
	// Iterate the proxies registered, and if the proxy using the given port is found, return it
	for _, proxy := range pm.proxies {
		if proxy.GetPort() == port && proxy.GetProtocol() == protocol {
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
	pe, err = pm.GetProxyFromParams(service.GetProtocol(), port)
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
