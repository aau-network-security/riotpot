package configuration

import (
	"os"
	"fmt"
	"log"
	"net"
	"sync"
	"bufio"
	"strings"

	"gorm.io/gorm"
	"github.com/riotpot/internal/greeting"
	"github.com/riotpot/pkg/services"
)

type Autopilot struct {
	Settings Settings
	Profile  Profile

	greeting greeting.Greet
	services services.Services
	wg       sync.WaitGroup
	DB       *gorm.DB

	loaded_plugins []string
	plugins_to_run []string
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

	// loads the services which are available for user to run
	a.loaded_plugins = a.services.GetServicesNames(a.services.GetServices())


	// set the plugins to run from the config file
	a.plugins_to_run = a.Settings.Riotpot.Start

	// check if the build is local or containerized
	if a.Settings.Riotpot.Local_build_on == "1" {
		// check if user want to run via config file or manually input
		// plugins to run.
		running_mode_decision := a.CheckRunningMode()

		// based on the user decision set the plugin running list
		if !running_mode_decision {
			// user decided to provide plugins manually
			a.plugins_to_run = a.GetPluginsFromUser()
		}
	}

	fmt.Printf("Plugins to run ")
	fmt.Println(a.plugins_to_run)

	// Check if the starting must be all the registered
	// or from the `Start` list.
	if a.Settings.Riotpot.Autod {
		a.services.RunAll()
	} else {
		for _, s := range a.plugins_to_run {
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

	service_paths := a.services.Autodiscover(a.Settings.Riotpot.Local_build_on)
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

// Converts the text separated by spaces into list items
func (a *Autopilot) TextToList(in string) (out []string) {
	out = strings.Fields(in)
	return out
}

// Reads the input from the terminal, returns the string
func (a *Autopilot) ReadInput() (text string) {
	reader := bufio.NewReader(os.Stdin)
	text, err := reader.ReadString('\n')
	if err != nil {
		log.Fatal(err)
	}

	return text
}

// Checks if the user wants to provide plugins to run manually
func (a *Autopilot) CheckRunningMode() (decision bool) {
	fmt.Print("Run plugins from configuation file? [y/n]")

	for {
		response := a.ReadInput()
		response = strings.ToLower(strings.TrimSpace(response))

		if response == "y" || response == "yes" {
			return true
		} else if response == "n" || response == "no" {
			return false
		} else{
			fmt.Printf("Please type Yes(y) or No(n) only\n")
		}
	}
}

// Validates if the plugins inputed by the users match the available plugins
// TODO: print all the invalid plugins not just the first one encountered
func (a *Autopilot) ValidatePlugin(input_plugins []string) (validated bool){
	for _, plugin := range input_plugins {
		validated := a.services.ValidatePluginByName(strings.Title(plugin), a.loaded_plugins)
		if !validated {
			fmt.Printf("\n[-] Entered plugin \"%s\" doesn't exist, please enter plugins again... \n", plugin)
			return false
		}
	}
	return true
}

// Gives which plugins user wants to load in RIoTPot
func (a *Autopilot) GetPluginsFromUser() (plugins []string) {
	for {
		fmt.Printf("\nPlugins available to run ")
		fmt.Println(a.loaded_plugins)
		fmt.Print("Enter the plugins to run separated by space: ")

		text := a.ReadInput()
		plugins = a.TextToList(text)
		validated := a.ValidatePlugin(plugins)
		if !validated {
			continue
		}
		break
	}
	return plugins
}
