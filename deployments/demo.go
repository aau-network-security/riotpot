package main

import (
	"fmt"
	"os"
	"log"
	"strings"
	// "gopkg.in/yaml.v3"
	// "github.com/gobuffalo/packr"
	errors "github.com/riotpot/tools/errors"
	"github.com/riotpot/tools/environ"
	"github.com/riotpot/tools/arrays"
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
	// box := packr.NewBox("../configs/samples")
	// data, err := box.Find("configuration.yml")

	// err = yaml.Unmarshal(data, &sett)

	// errors.Raise(err)
	a := configuration.Autopilot{
		Profile:  profile,
		Settings: sett,
	}

	a.Greeting()
	a.Settings.Riotpot.Start = arrays.StringToArray(a.Settings.Riotpot.Boot)
	fmt.Println(a.Settings.Riotpot.Start)
	a.RegisterPlugins()
	a.DiscoverImages()
	a.DiscoverRunningMode()
	a.SetLoadedPlugins()
	// fmt.Println(a.loaded_plugins)
	input_mode := a.CheckInteractionMode()
	existing_mode := a.Settings.Riotpot.Mode
	target_change := "s/mode: "+existing_mode+"/mode: "+input_mode+"/g"
	environ.ExecuteCmd("sed","-i", "-e", target_change, "../configs/samples/configuration.yml")

	if input_mode == "low" {
		fmt.Printf("Plugins available to run %q\n", a.Settings.Riotpot.Start)

		// user decided to provide plugins manually
		plugins_selected := a.GetPluginsFromUser()
		target_change = "s/boot: "+a.Settings.Riotpot.Boot+"/boot: "+strings.Join(plugins_selected, " ")+"/g"
		environ.ExecuteCmd("sed","-i", "-e", target_change, "../configs/samples/configuration.yml")
	} else if input_mode == "high" {
		fmt.Printf("\nDocker containers available to run ")
		fmt.Println(sett.GetDockerImages())
		fmt.Printf("\n")
		images := a.GetContainersFromUser()
		target_change = "s/start_images: "+a.Settings.Riotpot.Start_images+"/start_images: "+strings.Join(images, " ")+"/g"
		environ.ExecuteCmd("sed","-i", "-e", target_change, "../configs/samples/configuration.yml")
		FillConfig(images, &a)
	} else if input_mode == "hybrid" {
		fmt.Printf("Plugins available to run %q\n", a.Settings.Riotpot.Start)

		// user decided to provide plugins manually
		plugins_selected := a.GetPluginsFromUser()
		target_change = "s/boot: "+a.Settings.Riotpot.Boot+"/boot: "+strings.Join(plugins_selected, " ")+"/g"
		environ.ExecuteCmd("sed","-i", "-e", target_change, "../configs/samples/configuration.yml")

		fmt.Printf("\nDocker containers available to run ")
		fmt.Println(sett.GetDockerImages())
		fmt.Printf("\n")
		images := a.GetContainersFromUser()
		target_change = "s/start_images: "+a.Settings.Riotpot.Start_images+"/start_images: "+strings.Join(images, " ")+"/g"
		environ.ExecuteCmd("sed","-i", "-e", target_change, "../configs/samples/configuration.yml")
		FillConfig(images, &a)
	}
	
	fmt.Printf("Perfect!, now run the command 'docker-compose -f docker-compose.yml up -d --build'")
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


	// _, err = file.WriteString("The Go language was conceived in September 2007 by Robert Griesemer, Rob Pike, and Ken Thompson at Google.")
	if err != nil {
		log.Fatalf("failed writing to file: %s", err)
	}
}