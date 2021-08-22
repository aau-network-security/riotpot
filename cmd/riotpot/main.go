// Main Application
package main

import (
	"github.com/riotpot/tools/errors"
	"github.com/riotpot/internal/database"
	"github.com/riotpot/internal/configuration"
)

// `main` starts all the submodules containing the emulator services.
// It is the first function called when the application is run.
// It also acts as an orchestrator, which dictates the functioning of the application.
func main() {
	// Load the profile configuration
	profile, err := configuration.NewProfile()
	errors.Raise(err)

	// Load the configuration settings
	sett, err := configuration.NewSettings()
	errors.Raise(err)

	db := database.Database{
		// use only the default database
		Config: sett.Databases[0],
	}

	// load the connection to the database before anything
	conn := db.Connection()

	// For now, run with the autopilot on...
	auto := configuration.Autopilot{
		Profile:  profile,
		Settings: sett,
		DB:       conn,
	}

	auto.Start()
}
