/*
This package contains implementations of the services manager
*/
package services

import (
	"fmt"
	"path/filepath"
	"plugin"

	"github.com/riotpot/internal/globals"
	lr "github.com/riotpot/internal/logger"
	"github.com/riotpot/tools/errors"
	"golang.org/x/exp/slices"
)

var (
	pluginPath = "pkg/plugin/*/*.so"
	Services   = NewServiceManager(pluginPath)
)

// Function to get an stored service plugin.
// Note: the symbol used to get the plugin is "Name", which must be present in
// the plugin, and return type `Service` interface.
// based on: https://echorand.me/posts/getting-started-with-golang-plugins/
func getServicePlugin(path string) Service {

	// Open the plugin within the path
	pg, err := plugin.Open(path)
	errors.Raise(err)

	// check the name of the function that exports the service
	// The plugin *Must* contain a variable called `Plugin`.
	s, err := pg.Lookup("Plugin")
	errors.Raise(err)

	// log the name of the plugin being loaded
	fmt.Printf("Loading plugin: %s...\n", *s.(*string))

	// check if the reference symbol exists in the plugin
	rf, err := pg.Lookup(*s.(*string))
	errors.Raise(err)

	// Load the service in a variable as the interface Service.
	newservice := rf.(func() Service)()

	return newservice
}

// Get the plugin services included in the app
func pluginServices(pathLike string) (services []Service, err error) {
	// Get the paths to the plugins
	paths, err := filepath.Glob(pathLike)
	if err != nil {
		return
	}

	// Get the actual plugin and add it to the slice
	for _, path := range paths {
		service := getServicePlugin(path)
		services = append(services, service)
	}

	return
}

func RemovableService(service Service) (isRemovable bool) {
	// Add here the interfaces of services that should not be removable
	switch service.(type) {
	// Whether the service has the interface of a plugin
	case PluginService:
		isRemovable = true
	}

	return
}

type ServiceManager interface {
	// Register services
	addService(services ...Service) (serv []Service, err error)

	CreateService(name string, port int, network globals.Network, host string, interaction globals.Interaction) (Service, error)

	// Delete a service
	DeleteService(id string) (err error)

	// Get the list of services by their name
	GetServices() []Service

	// Get a single service
	GetService(id string) (Service, error)

	// Start the plugin services
	Start(ids ...string) ([]Service, error)
}

type ServiceManagerItem struct {
	ServiceManager

	// Set of services registered
	services []Service
}

// Add a service to the services map if it did not exist
func (se *ServiceManagerItem) addService(services ...Service) (serv []Service, err error) {
	// Returns a list of ID strings
	getServicesIDs := func(services []Service) (servs []string) {
		for _, service := range services {
			servs = append(servs, service.GetID())
		}
		return
	}

	// Convert the registered services into a simple array
	registeredIDs := getServicesIDs(se.GetServices())

	// Iterate the slice of services provided and add them to the services map
	for _, service := range services {
		serviceID := service.GetID()

		// Check whether the service is registered, and if not, add it to the list
		if !slices.Contains(registeredIDs, serviceID) {
			serv = append(serv, service)
		}
	}

	// Add the services to the registered services array
	se.services = append(se.services, serv...)

	return
}

// Creates a new service and register it in the manager
func (se *ServiceManagerItem) CreateService(name string, port int, network globals.Network, host string, interaction globals.Interaction) (s Service, err error) {
	// Iterate the services to determine whether the
	for _, service := range se.GetServices() {
		// Validate the name
		if service.GetName() == name {
			err = fmt.Errorf("service name already taken")
			return
		}

		// Validate the address
		if service.GetPort() == port && service.GetNetwork() == network && service.GetHost() == host {
			err = fmt.Errorf("service address already taken")
			return
		}
	}

	// Create the new service
	s = NewService(name, port, network, host, interaction)

	// Append the new service to the list
	se.services = append(se.services, s)
	return
}

// Remove a service from the list of registered
func (se *ServiceManagerItem) DeleteService(id string) (err error) {
	// Get all the services
	services := se.GetServices()

	for ind, service := range services {

		// Check if the service id is equal to the one received
		if service.GetID() == id && !service.IsLocked() {

			// Remove it from the slice by replacing it with the last item from the slice,
			// and reducing the slice by 1 element
			lastInd := len(services) - 1

			services[ind] = services[lastInd]
			se.services = services[:lastInd]

			return
		}
	}

	// If it was not found by this point, return an error
	err = fmt.Errorf("service not found")
	return
}

// Get services by name from the list of registered services
func (se *ServiceManagerItem) GetServices() []Service {
	return se.services
}

func (se *ServiceManagerItem) GetService(id string) (ret Service, err error) {
	for _, service := range se.GetServices() {
		if service.GetID() == id {
			ret = service
			return
		}
	}

	// If it was not found by this point, return an error
	err = fmt.Errorf("service not found")
	return
}

func (se *ServiceManagerItem) Start(ids ...string) (servs []Service, err error) {
	for _, id := range ids {
		service, e := se.GetService(id)
		if e != nil {
			err = e
			return
		}

		i, ok := service.(*PluginServiceItem)
		// If the service is not
		if !ok {
			err = fmt.Errorf("service %s is can not be started", service.GetName())
			return
		}

		// Run the service
		i.Run()
		servs = append(servs, i)
	}

	return
}

// Create a new pointer to a supervisor
func NewServiceManager(pluginPath string) (manager ServiceManager) {
	// Initialise the manager
	manager = &ServiceManagerItem{
		services: []Service{},
	}

	// Discover the services available to riotpot (running and stopped)
	services, err := pluginServices(pluginPath)
	if err != nil {
		lr.Log.Fatal().Err(err).Msgf("One or more services could not be found")
	}

	// Add/register the plugin services
	manager.addService(services...)

	return
}
