// This package implements the configuration of RiotPot.
// The package contains interfaces that help load, store and modify the configuration.
package configuration

import (
	"fmt"
	"os"
	"strconv"
	"strings"

	"github.com/riotpot/tools/arrays"
	environ "github.com/riotpot/tools/environ"
	errors "github.com/riotpot/tools/errors"

	"gopkg.in/yaml.v3"
)

func NewSettings() (s Settings, err error) {
	s = Settings{}
	err = s.Load("configs/samples/configuration.yml")
	s.ResolveEnv()
	return
}

// Interface that obligates the child to have the `load` and `save` methods
type Configuration interface {
	Load()
	Save()
}

// General configuration structure. It provides methods and attributes for parsing
// different types of configuration files, store, load and transform the state.
type Settings struct {
	Configuration

	Riotpot   ConfigRiotpot
	Databases []ConfigDatabase

	// Secret key string.
	Secret string
}

// Load the configuration on the child.
func (conf *Settings) Load(path string) (err error) {
	data, err := os.ReadFile(path)
	errors.Raise(err)

	err = yaml.Unmarshal(data, &conf)
	errors.Raise(err)

	return err
}

// Stores the configuration into the given path in `.yml` format.
func (conf *Settings) Save(path string) (err error) {
	// marshal the content of the configuration into a `.yaml` document
	d, err := yaml.Marshal(&conf)
	errors.Raise(err)

	// Save to file.
	// Mode 640: https://chmodcommand.com/chmod-640/
	// Note: this truncates the file if it already exists !!!
	err = os.WriteFile(path, d, 0640)
	errors.Raise(err)

	return err
}

// Validates the name of the emulator
func (conf *Settings) ValidateEmulators(service_paths []string) []string {
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
			fmt.Printf("[-] Plugin %s not allowed, skipping...\n", service)
		}
	}
	// Check if the array of emulators allowed contains the service
	return val
}

// This method overwrites the settings with the values from the environment
func (conf *Settings) ResolveEnv() {
	var err error

	// overwrite Autodiscover configuration setting
	if value, ok := os.LookupEnv("AUTOD"); ok {
		conf.Riotpot.Autod, err = strconv.ParseBool(value)
		errors.Raise(err)
	}

	// overwrite Starting emulators configuration setting
	if value, ok := os.LookupEnv("START"); ok {
		if value != "" {
			var emulators = strings.Split(value, ",")
			conf.Riotpot.Start = emulators
		}
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

	db_cfg.Identity.Name = environ.Getenv("DB_NAME", db_cfg.Identity.Name)

	conf.Databases[0] = db_cfg
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
	Autod bool
	// List of emulators that the application can access to.
	// This list will be evaluated against the `/emulators/` dir content.
	Emulators []string
	// List of emulators that must be run at start
	Start []string
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
