package main

import (
	"os"
	"fmt"
	"log"
	"strings"
	"github.com/riotpot/tools/arrays"
	"github.com/riotpot/tools/environ"
	errors "github.com/riotpot/tools/errors"
	"github.com/riotpot/internal/configuration"
)

func main() {
	// Resets the existing settings from previous runs
	environ.ExecuteCmd("cp", "../configs/samples/configuration-template.yml", "../configs/samples/configuration.yml")
	environ.ExecuteCmd("cp", "docker-compose-template.yml", "docker-compose.yml")

	profile, err := configuration.NewProfile()
	errors.Raise(err)

	sett, err := configuration.NewSettings()
	errors.Raise(err)

	a := configuration.Autopilot{
		Profile:  profile,
		Settings: sett,
	}

	a.Greeting()
	a.Settings.Riotpot.Start = arrays.StringToArray(a.Settings.Riotpot.Boot_plugins)

	supported_plugins := arrays.StringToArray(a.Settings.Riotpot.Boot_plugins)
	fmt.Printf("[+] Plugins available to run ")
	fmt.Println(supported_plugins)
	a.DiscoverImages()
	a.DiscoverRunningMode()
	a.SetLoadedPlugins()

	input_mode := a.CheckInteractionMode()
	existing_mode := a.Settings.Riotpot.Mode
	target_change := "s/mode: "+existing_mode+"/mode: "+input_mode+"/g"
	environ.ExecuteCmd("sed","-i", "-e", target_change, "../configs/samples/configuration.yml")

	if input_mode == "low" {
		fmt.Printf("Plugins available to run %q\n", a.Settings.Riotpot.Start)

		// user decided to provide plugins manually
		plugins_selected := a.GetPluginsFromUser()
		target_change = "s/boot_plugins: "+a.Settings.Riotpot.Boot_plugins+"/boot_plugins: "+strings.Join(plugins_selected, " ")+"/g"
		environ.ExecuteCmd("sed","-i", "-e", target_change, "../configs/samples/configuration.yml")
	} else if input_mode == "high" {
		fmt.Printf("\nDocker containers available to run are ")
		fmt.Println(sett.GetDockerImages())
		fmt.Printf("\n")
		images := a.GetContainersFromUser()
		target_change = "s/start_images: "+a.Settings.Riotpot.Start_images+"/start_images: "+strings.Join(images, " ")+"/g"
		environ.ExecuteCmd("sed","-i", "-e", target_change, "../configs/samples/configuration.yml")
		FillConfig(images, &a)
	} else if input_mode == "hybrid" {
		fmt.Printf("Plugins available to run are %q\n", a.Settings.Riotpot.Start)

		// user decided to provide plugins manually
		plugins_selected := a.GetPluginsFromUser()
		target_change = "s/boot_plugins: "+a.Settings.Riotpot.Boot_plugins+"/boot_plugins: "+strings.Join(plugins_selected, " ")+"/g"
		environ.ExecuteCmd("sed","-i", "-e", target_change, "../configs/samples/configuration.yml")

		fmt.Printf("\nDocker containers available to run are")
		fmt.Println(sett.GetDockerImages())
		fmt.Printf("\n")
		images := a.GetContainersFromUser()
		target_change = "s/start_images: "+a.Settings.Riotpot.Start_images+"/start_images: "+strings.Join(images, " ")+"/g"
		environ.ExecuteCmd("sed","-i", "-e", target_change, "../configs/samples/configuration.yml")
		FillConfig(images, &a)
	}
	// Ashi was here, Upon success following will be displayed
	fmt.Printf("Perfect!, now run the command\n") 
	fmt.Printf("\tdocker-compose -f docker-compose.yml up -d --build\n")
}

func FillConfig(images []string, a *configuration.Autopilot) {
	file, err := os.OpenFile("docker-compose.yml", os.O_WRONLY|os.O_APPEND, 0644)
	if err != nil {
		log.Fatalf("failed opening file: %s", err)
	}
	defer file.Close()

	for _,image := range images {
		_, err = file.WriteString("\n")
		_, err = file.WriteString("\n")
		image_tag := "  "+image+":"
		_, err = file.WriteString(image_tag)
		_, err = file.WriteString("\n")
		uri := a.Settings.GetContainerURI(image)
		image_option := "    image: "+uri
		_, err = file.WriteString(image_option)
		_, err = file.WriteString("\n")
		networks_tag := "    networks:"
		_, err = file.WriteString(networks_tag)
		_, err = file.WriteString("\n")
		_, err = file.WriteString("      honeypot:")
		_, err = file.WriteString("\n")
		ip := a.Settings.GetContainerIP(image)
		ip_addr_tag := "        ipv4_address: "+ip
		_, err = file.WriteString(ip_addr_tag)
	}

	if err != nil {
		log.Fatalf("failed writing to file: %s", err)
	}
}