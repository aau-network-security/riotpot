// This package implements the configuration of RiotPot.
// The package contains interfaces that help load, store and modify the configuration.
package configuration

import (
	"embed"
	"fmt"
	"log"
	"os"
	"strings"

	arrays "github.com/riotpot/tools/arrays"
	environ "github.com/riotpot/tools/environ"
	errors "github.com/riotpot/tools/errors"
	"gopkg.in/yaml.v3"
)

// Provides common identification attributes for each configuration.
// This structure must be private for each configuration object.
type IdentityConfiguration struct {
	ID   string
	Name string `yaml:"name"`
}

// RiotPot configuration structure. It includes all the attributes related to the riotpot framework.
// Moreover, it defines the emulators that must be loaded, and watches over them during runtime.
type RiotpotConfiguration struct {
	Identity *IdentityConfiguration
	// List of emulators that the application can access to.
	// This list will be evaluated against the `/pkg/plugin/` dir content.
	Emulators []string
	// List of plugins used for Riotpot to manage runtime plugins
	Start []string
}

// Database configuration structure. It gives an interface to load and access specific databases.
type DatabaseConfiguration struct {
	Identity *IdentityConfiguration

	// username to use in the db
	Username string
	// password for the user
	Password string
	// host in where the db is hosted
	Host string
	// port to connect to the database
	Port string
	// database name to connect
	Name string
}

// Interface that obligates the child to have the `load` and `save` methods
type ConfigurationProtocol interface {
	// Load settings method
	Load()
	// Save or store settings method
	Save()
}

// General configuration structure. It provides methods and attributes for parsing
// different types of configuration files, store, load and transform the state.
type Configuration struct {
	ConfigurationProtocol

	Riotpot  RiotpotConfiguration
	Database DatabaseConfiguration

	// Private fields
	// Configuration file path
	configPath string
}

// Load the configuration on the child.
func (conf *Configuration) Load() (err error) {
	// Read the content of the configuration file
	var f embed.FS
	data, _ := f.ReadFile(conf.configPath)

	// Serialise the content of the yaml file and load it into the structure
	err = yaml.Unmarshal(data, &conf)

	if err != nil {
		log.Fatal(err)
	}

	return err
}

// Stores the configuration into the given path in `.yml` format.
func (conf *Configuration) Save() (err error) {
	// marshal the content of the configuration into a `.yaml` document
	d, err := yaml.Marshal(&conf)
	errors.Raise(err)

	// Save to file.
	// Mode 640: https://chmodcommand.com/chmod-640/
	// Note: this truncates the file if it already exists !!!
	err = os.WriteFile(conf.configPath, d, 0640)
	errors.Raise(err)

	return err
}

// Validates the name of the emulator
func (conf *Configuration) ValidateEmulators(service_paths []string) []string {
	var val []string
	fmt.Printf("[+] Allowed plugins: %v\n", conf.Riotpot.Emulators)

	for _, p := range service_paths {
		//        '---path----'    ---plugin-----    -file-
		// split: `*pkg/plugin/` + `<plugin_name>/` + `*`
		sli := strings.Split(strings.SplitAfter(p, "pkg/plugin/")[1], "/")[0]
		// Transform the name of the plugin to lower case
		service := strings.ToLower(sli)

		// check if the service is in the allowed emulators slice.
		ok := arrays.Contains(conf.Riotpot.Emulators, service)
		if ok {
			val = append(val, p)
		} else {
			log.Printf("[-] Plugin %s not allowed, skipping...\n", service)
			conf.Riotpot.Start = arrays.DropItem(conf.Riotpot.Start, service)
		}
	}
	// Check if the array of emulators allowed contains the service
	return val
}

// This method overwrites the settings with the values from the environment
func (conf *Configuration) ResolveEnv() {

	// overwrite Starting emulators configuration setting
	if value, ok := os.LookupEnv("START"); ok {
		if value != "" {
			var emulators = strings.Split(value, ",")
			conf.Riotpot.Start = emulators
		}
	}

	// overwrite the default database to be used
	var dbCfg = conf.Database
	dbCfg = DatabaseConfiguration{
		Username: environ.Getenv("DB_USER", dbCfg.Username),
		Password: environ.Getenv("DB_PASS", dbCfg.Password),
		Host:     environ.Getenv("DB_HOST", dbCfg.Host),
		Port:     environ.Getenv("DB_PORT", dbCfg.Port),
		Name:     environ.Getenv("DB_NAME", dbCfg.Name),
	}

	dbCfg.Identity.Name = environ.Getenv("DB_NAME", dbCfg.Identity.Name)

	conf.Database = dbCfg
}

/**
Constructor for the configuration
This method loads the configuration from a local file
then, it overrides the configuration fields with the values from the environment variables
*/
func NewConfiguration() (conf Configuration, err error) {
	conf = Configuration{
		configPath: "configs/configuration.yml",
	}

	err = conf.Load()
	conf.ResolveEnv()
	return
}
