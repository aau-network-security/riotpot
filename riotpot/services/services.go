// This package provides multiple interfaces to load the services, validate them before running them
// and watching over their status
package services

import (
	"plugin"
	"riotpot/utils/errors"
)

// Function to get an stored service plugin.
// Note: the symbol used to get the plugin is "NewService", which must be present in
// the plugin, and return type `Service` interface.
// based on: https://echorand.me/posts/getting-started-with-golang-plugins/
func getServicePlugin(path string) Service {

	// Open the plugin within the path
	pg, err := plugin.Open(path)
	errors.Raise(err)

	// check if the reference symbol exists in the plugin
	rf, err := pg.Lookup("NewService")
	errors.Raise(err)

	// Load the service in a variable as the interface Service.
	newservice := rf.(func() Service)()

	return newservice
}

// Logger writes to the system log.
type Logger interface {
	Error(v ...interface{}) error
	Warning(v ...interface{}) error
	Info(v ...interface{}) error

	Errorf(format string, a ...interface{}) error
	Warningf(format string, a ...interface{}) error
	Infof(format string, a ...interface{}) error
}

// Interface used by every service plugin that offers a service. At the very least, every plugin
// must contain the set of methods and attributes from this interface.
// It is up to the plugin to determine the implementation of these methods for the most part.
type Service interface {
	Init(map[string]interface{})

	// The `run` method should be the implementation of starting the service.
	Run() error

	// The `stop` method should kill the process of the service gracefully.
	Stop() error

	// Tells the service to restart on premise.
	Restart() error

	// Returns the status of the service.
	Status() error

	// Interface to print the current logs of the service.
	Logger(errs chan<- error) (Logger, error)
}

// Wrapper for the individual services.
type Services struct {
	Id int

	// List of services registered the wrapper
	services []Service
}

// Method used to append a new service to the list of the wrapper
func (se *Services) Register(service Service) {
	se.services = append(se.services, service)
}

// This function utilizes a list of starting services
// to create and register new services.
//
//	Note: This function does not discern between new and already running services!
func (se *Services) AutoRegister(services []string) {
	// iterate through the slice of services
	for _, emu := range services {
		service := getServicePlugin(emu)
		se.Register(service)
	}
}
