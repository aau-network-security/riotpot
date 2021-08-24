package configuration

import (
	"os"
	"fmt"
	"log"
	"net"
	"sync"
	"bufio"
	"strconv"
	"strings"


	"github.com/riotpot/pkg/services"
	"github.com/riotpot/internal/greeting"
	"go.mongodb.org/mongo-driver/mongo"
	"github.com/riotpot/tools/arrays"
	"github.com/riotpot/tools/environ"
	"github.com/riotpot/pkg/profiles/ports"
)

type Autopilot struct {
	Settings Settings
	Profile  Profile

	greeting greeting.Greet
	services services.Services
	wg       sync.WaitGroup
	DB       *mongo.Client
	loaded_plugins []string
	loaded_containers []string
	plugins_to_run []string
	containers_to_run [] string
	interaction_mode string
	remote_host_ip string
	docker_context_name string
}

// Method to start the autopilot.
// It gets the list of emulators services available in the file system,
// and then it starts the given in the configuration file, given by either
// the `Autod` or `Start` variables.
func (a *Autopilot) Start() {
	// Say Hi, don't be rude!
	a.Greeting()

	a.wg = sync.WaitGroup{}
	a.wg.Add(1)
	a.Settings.Riotpot.Start = arrays.StringToArray(a.Settings.Riotpot.Boot_plugins)

	// register all the services plugins
	a.RegisterPlugins()
	a.DiscoverImages()

	// loads the services which are available for user to run
	a.SetLoadedPlugins()
	a.plugins_to_run = a.Settings.Riotpot.Start

	// check if the build is local or containerized
	if a.Settings.Riotpot.Local_build_on == "1" { // User is running the Riotpot in local build mode
		
		// check if user want to run via config file or manually input
		a.interaction_mode = a.CheckInteractionMode()

		if a.interaction_mode == "low" {
			// check if user want to run plugins via config file or manual input
			running_mode_decision := a.CheckRunningMode()

			if running_mode_decision == "manual" { // User has decided to provide plugins to run in interactive/manual way
				
				// user needs to choose which plugins to choose from this list 
				fmt.Printf("Plugins available to run ")
				fmt.Println(a.loaded_plugins)

				// user inputs the plugins to run
				a.plugins_to_run = a.GetPluginsFromUser()
				fmt.Println(a.plugins_to_run)
			} else {
				a.plugins_to_run = a.Settings.Riotpot.Start
				fmt.Printf("\nPlugins to run are ")
				fmt.Println(a.plugins_to_run)
				if !a.ValidatePlugin(a.plugins_to_run) {
					log.Fatalf("\nPlease check the config file, and try again\n")
				}
			}
		} else if a.interaction_mode == "high" {
			// reset the plugins since in high interaction mode local plugins are not to run
			a.plugins_to_run = nil
			// check if user wants to provide containers to run via config file or input them manually
			running_mode_decision := a.CheckRunningMode()

			if running_mode_decision == "manual" {
				fmt.Printf("\nDocker containers available to run ")
				fmt.Println(a.loaded_containers)
				fmt.Printf("\n")

				// User decided to provide containers manually, internally checks if the containers provided are valid or not
				a.containers_to_run = a.GetContainersFromUser()
				a.CheckContainerPort()
				a.DeployContainers()
			} else {
				a.containers_to_run = a.Settings.GetDockerImagesToRun()
				fmt.Printf("\nContainers to run are ")
				fmt.Println(a.containers_to_run)

				if !a.ValidateContainers(a.containers_to_run) {
					log.Fatalf("\nPlease check the config file, and try again\n")
				}

				// check if the port which containers require to run is free on host machine  
				a.CheckContainerPort()
				a.DeployContainers()
			}
		} else {
			// Hybrid mode
			running_mode_decision := a.CheckRunningMode()
			if running_mode_decision == "manual" {
				fmt.Printf("\nPlugins available to run ")
				fmt.Println(a.loaded_plugins)
				fmt.Printf("\n")
				a.plugins_to_run = a.GetPluginsFromUser()
				
				fmt.Printf("\nDocker containers available to run ")
				fmt.Println(a.loaded_containers)
				fmt.Printf("\n")
				a.containers_to_run = a.GetContainersFromUser()
				a.CheckContainerPort()
				a.DeployContainers()
			} else {
				a.plugins_to_run = a.Settings.Riotpot.Start
				a.containers_to_run = a.Settings.GetDockerImagesToRun()

				fmt.Printf("\nPlugins to run are ")
				fmt.Println(a.plugins_to_run)
				fmt.Printf("\nConainers to run are ")
				fmt.Println(a.containers_to_run)

				if !a.ValidatePlugin(a.plugins_to_run) {
					log.Fatalf("\nPlease check the config file, and try again\n")
				}
				if !a.ValidateContainers(a.containers_to_run) {
					log.Fatalf("\nPlease check the config file, and try again\n")
				}
				a.CheckContainerPort()
				a.DeployContainers()
			}
		}
	} else {
		// check which mode of run is set by user
		a.CheckModesFromConfig()

		if a.Settings.Riotpot.Mode == "low" {
			a.plugins_to_run = a.Settings.Riotpot.Start
			fmt.Printf("\nPlugins to run are ")
			fmt.Println(a.plugins_to_run)

			if !a.ValidatePlugin(a.plugins_to_run) {
					log.Fatalf("\nPlease check the config file, and try again\n")
			}
		} else if a.Settings.Riotpot.Mode == "high" {
			a.plugins_to_run = nil
			a.containers_to_run = a.Settings.GetDockerImagesToRun()
			fmt.Printf("\nContainers to run are ")
			fmt.Println(a.Settings.Riotpot.Start_images)

			if !a.ValidateContainers(a.containers_to_run) {
					log.Fatalf("\nPlease check the config file, and try again\n")
			}

			// glider forwards all the traffic on specific port to the respective service container
			a.DeployGlider()

		} else if a.Settings.Riotpot.Mode == "hybrid" {
			a.plugins_to_run = a.Settings.Riotpot.Start
			fmt.Printf("\nPlugins to run are ")
			fmt.Println(a.plugins_to_run)
			fmt.Printf("\nContainers to run are ")
			fmt.Println(a.Settings.Riotpot.Start_images)

			if !a.ValidatePlugin(a.plugins_to_run) {
					log.Fatalf("\nPlease check the config file, and try again\n")
			}

			if !a.ValidateContainers(a.containers_to_run) {
					log.Fatalf("\nPlease check the config file, and try again\n")
			}
			a.DeployGlider()
		}
	}

	// runs the Riotpot core
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

// check which interaction mode in supplied by the user from config file 
func (a *Autopilot) CheckModesFromConfig() {
	mode_received := arrays.StringToArray(a.Settings.Riotpot.Mode)

	if len(mode_received) > 1 {
		log.Fatalf("\nPlease enter only one mode in Riotpot config mode, i.e. low, high or hybrid\n")
	} else if len(mode_received) == 0 {
		log.Fatalf("\nPlease enter atleast one mode in Riotpot config mode, i.e. low, high or hybrid\n")
	}

	if ! arrays.Contains(a.Settings.Riotpot.Allowed_modes, mode_received[0]) {
		log.Fatalf("\n %q mode is invalid, only choose low, high or hybrid mode only in Riotpot config\n", mode_received[0])
	}
}

// Register the services plugins
func (a *Autopilot) RegisterPlugins() {
	a.services = services.Services{}

	service_paths := a.services.Autodiscover(a.Settings.Riotpot.Local_build_on)
	service_paths = a.Settings.ValidateEmulators(service_paths)

	a.services.AutoRegister(service_paths)
	a.services.AddDB(a.DB)
}

// Discover the available docker images
func (a *Autopilot) DiscoverImages() {
	a.loaded_containers = a.Settings.GetDockerImages()
	fmt.Printf("[+] Found %d docker images \n", len(a.loaded_containers))
	fmt.Printf("[+] Allowed Docker images ")
	fmt.Println(a.loaded_containers)
}

// Displays which is the current Riotpot running mode, i.e. low, high or hybrid
func (a *Autopilot) DiscoverRunningMode() {
	mode := "[+] Current mode of running is "+ a.Settings.Riotpot.Mode+"\n"
	fmt.Printf(mode)
}

// Load the greeting
func (a *Autopilot) Greeting() {
	a.greeting = greeting.Greet{
		Tutorial: a.Profile.Greet.Tutorial,
		Initial:  a.Profile.Greet.Initial,
	}

	a.greeting.Greeting()
}

// Reads the input from the terminal, returns string
func (a *Autopilot) ReadInput() (text string) {
	reader := bufio.NewReader(os.Stdin)
	text, err := reader.ReadString('\n')
	if err != nil {
		log.Fatal(err)
	}

	return text
}

// Check if the port of container to run is available for listening 
// TO-DO: print all the invalid containers not just the first one encountered
func (a *Autopilot) CheckContainerPort() {
	for _, container := range a.containers_to_run {
		port := ports.GetPort(arrays.AddSuffix(container, "d"))
		// Convert the Port from int to string
		if environ.CheckPortBusy(ports.GetProtocol(arrays.AddSuffix(container, "d")) , strconv.Itoa(port)) == false {
			log.Fatalf("[-] Port %d of Container %q is already busy on host, please free it first!", port, container)
		}
	}
}

// Deploy containers on host machine
func (a *Autopilot) DeployContainers() {
	for _, container := range a.containers_to_run {
		uri := a.Settings.GetContainerURI(container)
		// get the port number of a given container, currently all container must have an entry in ports file
		port := strconv.Itoa(ports.GetPort(arrays.AddSuffix(container, "d")))
		port_mapping := port+":"+port
		app := environ.GetPath("docker")
		environ.ExecuteBackgroundCmd(app, "run", "-p", port_mapping, uri )
		fmt.Printf("\nContainer %q, deployed \n", container)
	}
}

// Port forwarding using glider for container ports 
func (a *Autopilot) DeployGlider() {
	for _, container := range a.containers_to_run {
		port := strconv.Itoa(ports.GetPort(arrays.AddSuffix(container, "d")))
		protocol := ports.GetProtocol(arrays.AddSuffix(container, "d"))
		listener := protocol+"://:"+port
		remote_ip := a.Settings.GetContainerIP(container)
		forwarder := protocol+"://"+ remote_ip +":"+port
		fmt.Println(a.remote_host_ip)
		
		app := environ.GetPath("glider")
		
		environ.ExecuteBackgroundCmd(app, "-verbose", "-listen", listener, "-forward", forwarder)
	}
}

// Interactively checks if the user wants to provide plugins to run manually
func (a *Autopilot) CheckRunningMode() (string) {
	fmt.Print("Run plugins from configuation file? [y/n]")

	for {
		response := a.ReadInput()
		response = strings.ToLower(strings.TrimSpace(response))

		if response == "y" || response == "yes" {
			return "config"
		} else if response == "n" || response == "no" {
			return "manual"
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

// Validates if the plugins to run matches the available plugins
// TODO: print all the invalid plugins not just the first one encountered
func (a *Autopilot) ValidatePlugin(in_plugins []string) (validated bool){
	if arrays.HasDuplicateItems(in_plugins){
		fmt.Printf("\n[-] Entered plugins has duplicate entries, please enter again\n")
		return false
	}
	for _, plugin := range in_plugins {
		validated := arrays.Contains(a.loaded_plugins, strings.Title(strings.ToLower(plugin)))
		if !validated {
			fmt.Printf("\n[-] Entered plugin \"%s\" doesn't exist, please enter plugins again... \n", plugin)
			return false
		}
	}

	return true
}

// Validates if the containers to run matches the loaded containers
// TODO: print all the invalid containers not just the first one encountered
func (a *Autopilot) ValidateContainers(in_containers []string) (validated bool){
	if arrays.HasDuplicateItems(in_containers){
		fmt.Printf("\n[-] Entered containers has duplicate entries, please enter again\n")
		return false
	}

	for _, container := range in_containers {
		validated := arrays.Contains( a.loaded_containers, strings.ToLower(container))
		if !validated {
			fmt.Printf("\n[-] Entered container \"%s\" doesn't exist, please enter plugins again... \n", container)
			return false
		}
		contains := arrays.Contains(a.plugins_to_run, arrays.AddSuffix(container, "d"))
		if contains {
			fmt.Printf("\n[-] Entered container \"%s\" is already selected to run as a local plugin, please enter again\n", container)
			return false
		}
	}
	return true
}

// Interactively gets which plugins user wants to load in RIoTPot
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

// Interactively gets which container user wants to load in RIoTPot
func (a *Autopilot) GetContainersFromUser() (containers []string) {
	for {
		fmt.Print("Enter the containers to run separated by space: ")
		input := a.ReadInput()
		containers = arrays.StringToArray(input)
		validated := a.ValidateContainers(containers)

		if !validated {
			continue
		}
		break
	}

	return containers
}

// Gives which plugins user wants to load in RIoTPot
func (a *Autopilot) SetLoadedPlugins() {
	a.loaded_plugins = a.services.GetServicesNames(a.services.GetServices())
}

// Validates if the given docker context exists and if it is set to default
func (a *Autopilot) ValidateDefaultDockerContext(to_check string)  {
	path := environ.GetPath("docker")
	cmd_output	:= environ.ExecuteCmd(path, "context" , "ls")
	cmd_out_slice := arrays.StringToArray(cmd_output)
	val_position := arrays.GetItemPosition(cmd_out_slice, to_check)

	if val_position== -1 {
		log.Fatalf("Docker context %q, not found", to_check)
	}

	if cmd_out_slice[val_position+1] != "*" {
		log.Fatalf("Docker context %q, is not set to default", to_check)
	}
}
