/*
This package provides common settings necessary for the app
It makes the environment variables accessible, providing only the necessary
information and loading defaults.
*/
package settings

import (
	"fmt"
	"io/ioutil"
	"log"
	"strings"

	"riotpot/utils/arrays"
	"riotpot/utils/emulator"
	"riotpot/utils/environ"
)

// initialize the options from the environment
var AUTODISCOVER_EMULATORS string = environ.Getenv("AUTODISCOVER_EMULATORS", "")
var SERVICES string = environ.Getenv("SERVICES", "ALL")

var INSTALLED_EMULATORS = []string{
	"echod",
	"fakeshell",
	"httpd",
	"sshd",
	"telnetd",
}

// Init of the emulators in the app to be used
// The emulators will be the intersection of the services desired to be used
// and the instaled emulators
var VALIDATED_EMULATOR_LIST = validate_emulators_services(INSTALLED_EMULATORS, SERVICES)
var EMULATORS emulator.Emulators = *emulator.Register_emulators(VALIDATED_EMULATOR_LIST)

func get_installed_emulators() []string {

	// If the autodiscover variable is set, make a map of the current
	// directories under the "emulators" folder.
	if AUTODISCOVER_EMULATORS != "" {

		// get the list of directories in the emulators folder
		directories, err := ioutil.ReadDir("/emulators/")
		if err != nil {
			log.Fatal(err)
		}

		// append the name of the folder to the installed emulators
		for _, d := range directories {
			if d.IsDir() {
				INSTALLED_EMULATORS = append(INSTALLED_EMULATORS, d.Name())
			}

		}
	}

	return INSTALLED_EMULATORS
}

func validate_emulators_services(installed_emulators []string, services string) []string {
	/*
		Validate the services passed to the application against the installed emulators
		If any of the services defined is not included, it will throw a fatal.
	*/

	installed_emulators = get_installed_emulators()

	if len(installed_emulators) > 0 {
		// define a variable which will contain the intersection
		// from the services and the emulators
		var intersection []string

		if services == "ALL" {
			intersection = installed_emulators
		} else {
			// spit the services string from the environment by comma into an array
			services := strings.Split(services, ",")

			// iterate through the services to check if they exists in the installed
			// emulators; otherwise throw an error
			for _, service := range services {
				if arrays.Contains(installed_emulators, service) {
					intersection = append(intersection, service)
				} else {
					// crash the loop and do not recover.
					errorString := fmt.Sprintf("runtime error: service %s was not found in the installed emulators", service)
					panic(errorString)
				}
			}
		}

		return intersection
	} else {
		panic("runtime error: there are no installed emulators")
	}
}
