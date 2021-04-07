package configuration

import (
	"encoding/json"
	"fmt"
	"os"
	"testing"
)

var (
	settings = Settings{}
	path     = "."
	fLoad    = "configuration"
	fSave    = "temp"
	err      error
)

// Test loading and saving configuration from a file to another
func TestLoadAndSaveConf(t *testing.T) {

	// load the configuration
	err = settings.Load(path, fLoad)
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
	err = settings.Save(path, fSave)
	if err != nil {
		t.Error(err)
	}

	// read the file just saved and prints it
	fp := fmt.Sprintf("%s/%s.yml", path, fSave)
	data, err := os.ReadFile(fp)
	if err != nil {
		t.Error(err)
	}

	t.Log(string(data))
}
