package configuration

import (
	"os"
	"fmt"
	"log"
	"net"
	"sync"
	"bufio"
	"strings"
	// "os/exec"

	"gorm.io/gorm"
	"github.com/riotpot/pkg/services"
	
	"github.com/riotpot/internal/greeting"
	"github.com/riotpot/tools/environ"
	"github.com/riotpot/tools/arrays"
)

type Autopilot struct {
	Settings Settings
	Profile  Profile

	greeting greeting.Greet
	services services.Services
	wg       sync.WaitGroup
	DB       *gorm.DB

	loaded_plugins []string
	loaded_containers []string
	plugins_to_run []string
	contianers_to_run [] string
	interaction_mode string
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

	// register all the services plugins
	a.RegisterPlugins()
	a.DiscoverImages()
	_ = environ.CheckDockerExists("mongodb")

	// loads the services which are available for user to run
	a.loaded_plugins = a.services.GetServicesNames(a.services.GetServices())

	// set the plugins to run from the config file
	a.plugins_to_run = a.Settings.Riotpot.Start

	// check if the build is local or containerized
	if a.Settings.Riotpot.Local_build_on == "1" {
		a.interaction_mode = a.CheckInteractionMode()

		if a.interaction_mode == "low" {
			// check if user want to run via config file or manually input
			// plugins to run.
			// _ = a.CheckInteractionMode()
			running_mode_decision := a.CheckRunningMode()

			// based on the user decision set the plugin running list
			if !running_mode_decision {
				fmt.Printf("Plugins available to run ")
				fmt.Println(a.loaded_plugins)

				// user decided to provide plugins manually
				a.plugins_to_run = a.GetPluginsFromUser()
			}	
		} else if a.interaction_mode == "high" {
			fmt.Printf("Docker images available to run ")
			fmt.Println(a.loaded_containers)

			running_mode_decision := a.CheckContainersRunMode()
			if !running_mode_decision {
				fmt.Printf("Containers available to run ")
				fmt.Println(a.Settings.GetDockerImages())

				// user decided to provide contianers manually
				a.contianers_to_run = a.GetContainersFromUser()
			}


		} else {
			fmt.Printf("\nPlugins available to run ")
			fmt.Println(a.loaded_plugins)
			fmt.Printf("\n")
			fmt.Printf("Docker images available to run ")
			fmt.Println(a.loaded_containers)
			fmt.Printf("\n")

		}
	}

	// fmt.Printf("Plugins to run ")
	// fmt.Println(a.plugins_to_run)

	// // cmd, err := exec.Command("/bin/sh", "bash.sh").Output()
	// // cmd, err := exec.LookPath("go")
	// // if err != nil {
	// //     log.Fatalf("[!] Error: %v", err)
 // //    }

 // //    output := string(cmd)
	// // fmt.Println(output)

	// server := environ.CheckPortBusy("tcp", ":22")
	// if server != true {
	// 	fmt.Println("Port is busy")
	// }
	// fmt.Println(server)	

	// _, exists := environ.GetPath("glider")
	// if exists {
	// 	a.ValidateDefaultDockerContext("docker-test")
	// 	_ 	= environ.ExecuteCmd("go", "version")	
	// } 

	// fmt.Println(demo1)		

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

// Register the services plugins
func (a *Autopilot) RegisterPlugins() {
	a.services = services.Services{}

	service_paths := a.services.Autodiscover(a.Settings.Riotpot.Local_build_on)
	service_paths = a.Settings.ValidateEmulators(service_paths)

	a.services.AutoRegister(service_paths)

	a.services.AddDB(a.DB)
}

// Discover the docker images available
func (a *Autopilot) DiscoverImages() {
	a.loaded_containers = a.Settings.GetDockerImages()
	fmt.Printf("[+] Found %d docker images \n", len(a.loaded_containers))
	fmt.Printf("[+] Allowed Docker images ")
	fmt.Println(a.loaded_containers)

}

// Load the greeting
func (a *Autopilot) Greeting() {
	a.greeting = greeting.Greet{
		Tutorial: a.Profile.Greet.Tutorial,
		Initial:  a.Profile.Greet.Initial,
	}

	a.greeting.Greeting()
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

// Checks if the user wants to run containers manually
func (a *Autopilot) CheckContainersRunMode() (decision bool) {
	fmt.Print("Run containers from configuation file? [y/n]")

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

// Checks in which mode user wants to run RIoTPot
func (a *Autopilot) CheckInteractionMode() (decision string) {
	fmt.Printf("\nSelect RIoTPot mode, Low-interaction mode, High-interaction mode or Hybrid-mode? [l,h,hy] [low, high, hybrid] \n")

	for {
		response := a.ReadInput()
		response = strings.ToLower(strings.TrimSpace(response))

		if response == "l" || response == "low" {
			return "low"
		} else if response == "h" || response == "high" {
			return "high"
		} else if response == "hy" || response == "hybrid" {
			return "hybrid"
		} else{
			fmt.Printf("Please type low(l) or high(h) or hybrid(hy) only\n")
		}
	}
}

// Validates if the plugins inputed by the user matches the available plugins
// TODO: print all the invalid plugins not just the first one encountered
func (a *Autopilot) ValidatePlugin(input_plugins []string) (validated bool){
	for _, plugin := range input_plugins {
		validated := arrays.Contains(a.loaded_plugins, strings.Title(strings.ToLower(plugin)))
		if !validated {
			fmt.Printf("\n[-] Entered plugin \"%s\" doesn't exist, please enter plugins again... \n", plugin)
			return false
		}
	}
	return true
}

// Validates if the contianers inputed by the users match the loaded containers
// TODO: print all the invalid containers not just the first one encountered
func (a *Autopilot) ValidateContianers(input_containers []string) (validated bool){
	for _, container := range input_containers {
		validated := arrays.Contains( a.loaded_containers, strings.ToLower(container))
		if !validated {
			fmt.Printf("\n[-] Entered container \"%s\" doesn't exist, please enter plugins again... \n", container)
			return false
		}
	}
	return true
}

// Gives which plugins user wants to load in RIoTPot
func (a *Autopilot) GetPluginsFromUser() (plugins []string) {
	for {
		fmt.Print("Enter the plugins to run separated by space: ")

		input := a.ReadInput()
		plugins = arrays.StringToArray(input)

		validated := a.ValidatePlugin(plugins)
		if !validated {
			continue
		}
		break
	}
	return plugins
}

// Gives which plugins user wants to load in RIoTPot
func (a *Autopilot) GetContainersFromUser() (containers []string) {
	for {
		fmt.Print("Enter the contianers to run separated by space: ")

		input := a.ReadInput()
		containers = arrays.StringToArray(input)
		validated := a.ValidateContianers(containers)

		if !validated {
			continue
		}
		break
	}

	return containers
}


// Validates if the given docker context exists and set to default
func (a *Autopilot) ValidateDefaultDockerContext(to_check string)  {
	cmd_output	:= environ.ExecuteCmd("docker", "context" , "ls")
	cmd_out_slice := arrays.StringToArray(cmd_output)
	val_position := arrays.GetItemPosition(cmd_out_slice, to_check)

	if val_position== -1 {
		log.Fatalf("Docker context %q, not found", to_check)
	}
	if cmd_out_slice[val_position+1] != "*" {
		log.Fatalf("Docker context %q, is not set to default", to_check)
	}
}
