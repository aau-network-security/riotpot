package configuration

import (
	"encoding/json"
	"os"
	"testing"

	cfg "github.com/riotpot/internal/configuration"
)

var (
	settings = cfg.Settings{}
	path     = "./configuration_test.yml"
)

// Test loading and saving configuration from a file to another
func TestLoadAndSaveConf(t *testing.T) {
	// load the configuration
	err := settings.Load(path)
	if err != nil {
		t.Error(err)
	}

	// pretty print marshalling to json
	out, err := json.MarshalIndent(settings, "", "	")
	if err != nil {
		t.Error(err)
	}

	t.Logf("%s\n", string(out))

	// save the configuration
	err = settings.Save(path)
	if err != nil {
		t.Error(err)
	}

	// read the file just saved and prints it
	data, err := os.ReadFile(path)
	if err != nil {
		t.Error(err)
	}

	t.Log(string(data))

}
