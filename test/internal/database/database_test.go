package database

import (
	"fmt"
	"testing"

	"github.com/riotpot/internal/configuration"
	"github.com/riotpot/internal/database"

	"gorm.io/gorm"
)

type UserTest struct {
	gorm.Model
	Name string
}

// Runs a simple test to check the health of a database and the configuration.
// Keep in mind this function requires the role `superuser` and the `superuser` db
// to previously exists!
// Example:
//  $ # create the user as a superuser and the db
//  $ createuser superuser -s
//  $ createdb superuser
func TestDatabase(t *testing.T) {

	var (
		// Load an identity for the database
		id = configuration.ConfigIdentity{
			Name: "test_db",
		}

		// Load a configuration for the database
		config = configuration.ConfigDatabase{
			Username: "superuser",
			Password: "",
			Host:     "127.0.0.1",
			Identity: id,
			Port:     "5432",
		}

		// Load a database object
		db = database.Database{
			Config: config,
		}

		// Load a random user
		user = UserTest{
			Name: "Test",
		}
	)

	// Connect to the db
	conn := db.Connection()

	// create and move to the db
	conn = conn.Exec(fmt.Sprintf("CREATE DATABASE %s", db.Config.Identity.Name))

	// migrate the model
	conn.AutoMigrate(&UserTest{})

	// Insert the user in the database
	result := conn.Create(&user)
	if result.Error != nil {
		t.Error(result.Error)
	}

	// Log some of the results...
	t.Logf("User ID: %v", user.ID)
	t.Logf("User Name: %v", user.Name)
}
