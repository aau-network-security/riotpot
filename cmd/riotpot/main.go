// Main Application
package main

import (
	"github.com/riotpot/internal/configuration"
	"github.com/riotpot/internal/database"
	"github.com/riotpot/tools/errors"
)

// `main` starts all the submodules containing the emulator services.
// It is the first function called when the application is run.
// It also acts as an orchestrator, which dictates the functioning of the application.
func main() {
	// Load the configuration settings
	conf, err := configuration.NewConfiguration()
	errors.Raise(err)

	db := database.Database{
		// use only the default database
		Config: conf.Databases[0],
	}

	// load the connection to the database before anything
	conn := db.Connection()

	// For now, run with the autopilot on...
	auto := configuration.Autopilot{
		Configuration: conf,
		DB:            conn,
	}

	auto.Start()
}
