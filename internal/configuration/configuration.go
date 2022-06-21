// This package implements the configuration of RiotPot.
// The package contains interfaces that help load, store and modify the configuration.
package configuration

import (
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"

	"github.com/gobuffalo/packr"
	arrays "github.com/riotpot/tools/arrays"
	environ "github.com/riotpot/tools/environ"
	errors "github.com/riotpot/tools/errors"
	"gopkg.in/yaml.v3"
)

/** Constructor for the configuration
*	This method loads the configuration from a local file
*	then, it overrides the configuration fields with the values from the environment variables
 */
func NewConfiguration() (conf Configuration, err error) {
	conf = Configuration{}
	conf.configPath = "configuration.yml"

	err = conf.Load()
	conf.ResolveEnv()
	return
}

// Interface that obligates the child to have the `load` and `save` methods
type ConfigurationInterface interface {
	// Load settings method
	Load()
	// Save or store settings method
	Save()
}

// General configuration structure. It provides methods and attributes for parsing
// different types of configuration files, store, load and transform the state.
type Configuration struct {
	ConfigurationInterface

	Riotpot   ConfigRiotpot //TODO This could be just a manager
	Databases []ConfigDatabase

	// Private fields
	// Configuration file path
	configPath string
}

// Load the configuration on the child.
func (conf *Configuration) Load() (err error) {
	// TODO: Why are we user packer to find a file in the system?
	box := packr.NewBox("../../configs/samples")
	data, _ := box.Find(conf.configPath)

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

// Retrieve the image name from Images tag in configuration file
// TODO: Delete this function. If the attribute needs to be an array, set it as an array
func (conf *Configuration) GetDockerImages() (images []string) {
	for _, val := range conf.Riotpot.Images {
		images = append(images, strings.TrimSuffix(arrays.StringToArray(val)[0], ","))
	}

	return images
}

// Retrieve the image name from Start_images tag in configuration file
// TODO: Delete this function. If the attribute needs to be an array, set it as an array
func (conf *Configuration) GetDockerImagesToRun() []string {
	// for _, val := range conf.Riotpot.Start_images {
	// 	images = append(images, strings.TrimSuffix(arrays.StringToArray(val)[0], ","))
	// }

	return arrays.StringToArray(conf.Riotpot.Start_images)
}

// Retrieve the container uri from Images tag in configuration file
// TODO: instead of assigning IP's, use container names
func (conf *Configuration) GetContainerURI(container string) (uri string) {
	for _, val := range conf.Riotpot.Images {
		data := strings.Split(val, ",")
		service := data[0]
		if container == service {
			uri = strings.TrimSpace(data[1])
			break
		}
	}

	return uri
}

// Retrieve the container image IP from Images tag in configuration file
// TODO: instead of assigning IP's, use container names
func (conf *Configuration) GetContainerIP(container string) (ip string) {
	for _, val := range conf.Riotpot.Images {
		data := strings.Split(val, ",")
		image := data[0]
		if container == image {
			ip = strings.TrimSpace(data[2])
			break
		}
	}

	return ip
}

// Retrieve the loaded plugins, i.e. plugins which are loaded in the system
// TODO: What is the purpose of this function? can't we access this from the class?
func (conf *Configuration) GetLoadedPlugins() (plugins []string) {
	return conf.Riotpot.Start
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
			fmt.Printf("[-] Plugin %s not allowed, skipping...\n", service)
			conf.Riotpot.Start = arrays.DropItem(conf.Riotpot.Start, service)
		}
	}
	// Check if the array of emulators allowed contains the service
	return val
}

// This method overwrites the settings with the values from the environment
func (conf *Configuration) ResolveEnv() {
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
		Dbname:   environ.Getenv("DB_NAME", db_cfg.Dbname),
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
	// List of plugins used for Riotpot to manage runtime plugins
	Start []string
	// Plugins which are booted in the system to run, supplied by user
	// TODO: Why is this a string and not an array? Also, isn't this the "start" attribute??
	Boot_plugins string
	// Variable to check if the run is for local build or not
	// TODO: Why is this a string and not a boolean?
	Local_build_on string
	// Available docker images along with docker registry name and ip address(for contianerized runs)
	Images []string
	// Interaction mode of Riotpot, used in containazried build
	// TODO: Change this from a string to an integer or a tuple
	Mode string
	// Modes of operation which are currently supported by Riotpot
	// TODO: Change this from a list of strings to a class instance. If there are modes, these should be part of the app
	Allowed_modes []string
	// Container images which are finalized to run in the Riotpot run
	// TODO: Why is this a string?
	Start_images string
}

// Database configuration structure. It gives an interface to load and access specific databases.
type ConfigDatabase struct {
	Identity ConfigIdentity

	/* Database configuration */
	// engine used in the db e.g. sql, postgres, sqlite
	Engine string // TODO: it seems that the VIP has decided that this is a MongoDB
	// username to use in the db
	Username string
	// password for the user
	Password string
	// host in where the db is hosted
	Host string
	// port to connect to the database
	Port string
	// database name to connect
	Dbname string
}
