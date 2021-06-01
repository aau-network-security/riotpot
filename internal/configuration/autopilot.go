package configuration

import (
	"fmt"
	"net"
	"sync"

	"github.com/riotpot/internal/greeting"
	"github.com/riotpot/pkg/services"
	"gorm.io/gorm"
)

type Autopilot struct {
	Settings Settings
	Profile  Profile

	greeting greeting.Greet
	services services.Services
	wg       sync.WaitGroup
	DB       *gorm.DB
}

// Method to start the autopilot.
// It gets the list of emulators services available in the file system,
// and then it starts the given in the configuration file, given by either
// the `Autod` or `Start` variables.
func (a *Autopilot) Start() {
	// Say Hi, don't be rude!
	a.Greeting()

	// block the main thread until we kill all the services
	a.wg = sync.WaitGroup{}
	a.wg.Add(1)

	// register all the services
	a.RegisterServices()

	// Check if the starting must be all the registered
	// or from the `Start` list.
	if a.Settings.Riotpot.Autod {
		a.services.RunAll()
	} else {
		for _, s := range a.Settings.Riotpot.Start {
			// get the service and run it
			service := a.services.Get(s)
			fmt.Printf("Starting %s...\n", service.GetName())

			if a.available(service.GetName(), service.GetPort()) {
				go service.Run()
			}
		}
	}

	a.wg.Wait()
}

func (a *Autopilot) available(name string, port int) (available bool) {
	// convert the port to a string
	sport := fmt.Sprintf("%d", port)
	ln, err := net.Listen("tcp", ":"+sport)

	if err != nil {
		fmt.Printf("Port %s unavailable, skipping %s ...\n", sport, name)
		return false
	}

	_ = ln.Close()
	return true
}

// Register the services
func (a *Autopilot) RegisterServices() {
	a.services = services.Services{}

	service_paths := a.services.Autodiscover()
	service_paths = a.Settings.ValidateEmulators(service_paths)

	a.services.AutoRegister(service_paths)
	a.services.AddDB(a.DB)
}

// Load the greeting
func (a *Autopilot) Greeting() {
	a.greeting = greeting.Greet{
		Tutorial: a.Profile.Greet.Tutorial,
		Initial:  a.Profile.Greet.Initial,
	}

	a.greeting.Greeting()
}
