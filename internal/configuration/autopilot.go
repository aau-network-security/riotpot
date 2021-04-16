package configuration

import (
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
	a.greeting.Greeting()

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
			go service.Run()
		}
	}

	a.wg.Wait()
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
	}
}
