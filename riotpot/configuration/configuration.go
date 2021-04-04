// This package implements the configuration of RiotPot.
// The package contains interfaces that help load, store and modify the configuration.
package configuration

import (
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"

	arrays "riotpot/utils/arrays"
	environ "riotpot/utils/environ"
	errors "riotpot/utils/errors"

	"gopkg.in/yaml.v3"
)

// General configuration structure. It provides methods and attributes for parsing
// different types of configuration files, store, load and transform the state.
type Settings struct {
	Riotpot   ConfigRiotpot
	Databases []ConfigDatabase

	// Secret key string.
	secret string
}

// Takes configuration data as string and updates the state of a given configuration based on it.
func (conf *Settings) Load(path string, filename string) error {
	// full path of the file.
	var filepath = fmt.Sprintf("%s/%s.yml", path, filename)

	// load the data into the memory
	data, err := os.ReadFile(filepath)
	errors.Raise(err)

	// unmarshal the data into the configuration settings.
	err = yaml.Unmarshal(data, &conf)
	errors.Raise(err)

	return err
}

// Stores the configuration into the given path in `.yml` format.
func (conf Settings) Save(path string, filename string) error {
	// marshal the content of the configuration into a `.yaml` document
	d, err := yaml.Marshal(&conf)
	errors.Raise(err)

	// full path of the file.
	var filepath = fmt.Sprintf("%s/%s.yml", path, filename)

	// Save to file.
	// Mode 640: https://chmodcommand.com/chmod-640/
	// Note: this truncates the file if it already exists !!!
	err = os.WriteFile(filepath, d, 0640)
	errors.Raise(err)

	return err
}

// This method overwrites the settings with the values from the environment
func (conf Settings) ResolveEnv() {
	var err error

	// overwrite Autodiscover configuration setting
	if value, ok := os.LookupEnv("AUTODISCOVER"); ok {
		conf.Riotpot.Autodiscover, err = strconv.ParseBool(value)
		errors.Raise(err)
	}

	// overwrite Starting emulators configuration setting
	if value, ok := os.LookupEnv("START"); ok {
		var emulators = strings.Split(value, ",")
		conf.Riotpot.Emulators = emulators
	}

	// overwrite the default database to be used
	var db_cfg = conf.Databases[0]
	db_cfg = ConfigDatabase{
		Engine:   environ.Getenv("DB_ENGINE", db_cfg.Engine),
		Username: environ.Getenv("DB_USER", db_cfg.Username),
		Password: environ.Getenv("DB_PASS", db_cfg.Password),
		Host:     environ.Getenv("DB_HOST", db_cfg.Host),
		Port:     environ.Getenv("DB_PORT", db_cfg.Port),
	}

}

// Provides common identification attributes for each configuration.
// This structure must be private for each configuration object.
type ConfigIdentity struct {
	ID   string
	Name string `yaml:"name"`
}

// RiotPot configuration structure. It includes all the attributes related to the riotpot framework.
// Moreover, it defines the emulators that must be loaded, and watches over them during runtime.
type ConfigRiotpot struct {
	Identity ConfigIdentity

	/* Riotpot configuration attributes: */
	// Defines if the emulators must be loaded directly from the file system.
	Autodiscover bool
	// List of emulators that the application can access to.
	// This list will be evaluated against the `/emulators/` dir content.
	Emulators []string
	// List of emulators that must be run at start
	Start []string
}

// Method used to verify the validity of the emulators.
// On success, it returns the list of emulators.
func (conf *ConfigRiotpot) Validated() []string {
	var directories []string

	// get the list of directories in the emulators folder
	documents, err := os.ReadDir("/emulators/")
	errors.Raise(err)

	// iterate through the documents stored in the directory.
	for _, d := range documents {
		if d.IsDir() {
			directories = append(directories, d.Name())
		}
	}

	// iterate through the list of emulators given as argument
	for _, emu := range conf.Emulators {
		if !arrays.Contains(directories, emu) {
			log.Fatalf("Error: emulator %s not found", emu)
		}
	}

	return conf.Emulators
}

// Database configuration structure. It gives an interface to load and access specific databases.
type ConfigDatabase struct {
	Identity ConfigIdentity

	/* Database configuration */
	// engine used in the db e.g. sql, postgres, sqlite
	Engine string
	// username to use in the db
	Username string
	// password for the user
	Password string
	// host in where the db is hosted
	Host string
	// port to connect to the database
	Port string
}
