// This package provides multiple interfaces to load the services, validate them before running them
// and watching over their status
package services

import (
	"fmt"
	"path/filepath"
	"plugin"
	"strings"

	"github.com/riotpot/tools/errors"
	"gorm.io/gorm"
)

// Function to get an stored service plugin.
// Note: the symbol used to get the plugin is "NewService", which must be present in
// the plugin, and return type `Service` interface.
// based on: https://echorand.me/posts/getting-started-with-golang-plugins/
func getServicePlugin(path string) Service {

	// Open the plugin within the path
	pg, err := plugin.Open(path)
	errors.Raise(err)

	// check the name of the function that exports the service
	// The plugin *Must* contain a variable called `Name`.
	s, err := pg.Lookup("Name")
	errors.Raise(err)

	// log the name of the plugin being loaded
	fmt.Printf("[+] Loading plugin: %s...\n", *s.(*string))

	// check if the reference symbol exists in the plugin
	rf, err := pg.Lookup(*s.(*string))
	errors.Raise(err)

	// Load the service in a variable as the interface Service.
	newservice := rf.(func() Service)()

	return newservice
}

// Interface used by every service plugin that offers a service. At the very least, every plugin
// must contain the set of methods and attributes from this interface.
// It is up to the plugin to determine the implementation of these methods for the most part.
type Service interface {
	Init(map[string]interface{})

	// Gets the name of the service
	GetName() string
	GetProtocol() string
	GetPort() int

	// The `run` method should be the implementation of starting the service.
	Run() error

	// The `stop` method should kill the process of the service gracefully.
	Stop() error

	// Tells the service to restart on premise.
	Restart() error

	// Returns the status of the service.
	Status() string

	// Set the database connection
	SetDb(conn *gorm.DB)
}

// Implements a mixin service that can be used as a base for any other service `struct` type.
type MixinService struct {
	// require the methods described by `Service` on loading
	Service

	// A connection to the database that must be initialized
	conn *gorm.DB

	// it is recommended to include some kind of identity
	// for the service.
	Id   int
	Name string

	// declare here any other variable relevant for the
	// service to run. This are merely examples...
	Protocol string
	Port     int
	Host     string

	// Stopping channel, a signal for the program to stop.
	StopCh chan int

	// boolean indicating if the service is running
	Running chan bool
}

func (mx *MixinService) Stop() error {
	var err error

	// checks if the channel is open
	isOpen := func(ch <-chan int) bool {
		select {
		case <-ch:
			// return true if we can read from the channel
			return true
		default:
		}

		return false
	}

	// send a stop signal to the channel so we can stop all the current
	// connections gracefully if the channel is open
	if isOpen(mx.StopCh) {
		mx.StopCh <- 1
	} else {
		err = fmt.Errorf("Service %s currently not running", mx.Name)
	}

	return err
}

func (mx *MixinService) Restart() (err error) {

	// Stops the service and calls `Run` again on it.
	// TODO: add storing of the current status of the service
	mx.Stop()
	go mx.Run()

	return err
}

func (mx *MixinService) GetName() string {
	return mx.Name
}

func (mx *MixinService) GetProtocol() string {
	return mx.Protocol
}

func (mx *MixinService) GetPort() int {
	return mx.Port
}

// Simple function on the mixin that checks if the service is
// currently running.
func (mx *MixinService) Status() string {
	// with select we prevent locking the thread
	select {
	case running := <-mx.Running:
		if running {
			return "Running"
		} else {
			return "Stopped"
		}
	default:
		return "Stopped"
	}
}

func (mx *MixinService) SetDb(conn *gorm.DB) {
	mx.conn = conn
}

func (mx *MixinService) Migrate(model interface{}) {
	if mx.conn != nil {
		mx.conn.AutoMigrate(model)
	} else {
		fmt.Print("Database not accessible")
	}
}

func (mx *MixinService) Store(model interface{}) {
	if mx.conn != nil {
		mx.conn.Create(model)
	} else {
		fmt.Print("Database not accessible")
	}
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
func (se *Services) AutoRegister(service_paths []string) {
	// iterate through the slice of services
	for _, emu := range service_paths {
		service := getServicePlugin(emu)
		se.Register(service)
	}
}

// Get a service by the name from the list of registered services
func (se *Services) Get(name string) (s Service) {
	// capitalize the name, so it matches the name, as it must be
	// an exported function
	name = strings.Title(name)

	// Iterate through the services to get the one with the name.
	// This method is rather slow, but the services won't normally
	// exceed 10 or 20.
	for _, service := range se.services {
		if n := service.GetName(); n == name {
			s = service
			return
		}
	}

	fmt.Printf("service not found: %s", name)
	return
}

// Get multiple registered services given their name
func (se *Services) GetMultiple(service_names ...string) (services []Service) {
	for _, service_name := range service_names {
		tempServ := se.Get(service_name)
		services = append(services, tempServ)
	}
	return
}

// Method to register all the services in the `pkg/plugin` folder
// The method looks for the `.so` plugin inside of the folders,
// doesn't matter the name of the folder.
func (se *Services) Autodiscover() []string {
	all, err := filepath.Glob("pkg/plugin/*/*.so")
	if err != nil {
		panic(err)
	}

	fmt.Printf("[+] Found %d plugins\n", len(all))
	return all
}

// Simple function that iterates through the registered services and
// starts them using their `Run` method.
// The method does not accept any argument.
func (se *Services) RunAll() {
	for _, s := range se.services {
		go s.Run()
	}
}

// Add a database connection to all the services
func (se *Services) AddDB(conn *gorm.DB) {
	for _, s := range se.services {
		s.SetDb(conn)
	}
}
