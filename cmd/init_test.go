package cmd

import (
	"encoding/json"
	"testing"

	"github.com/spf13/afero"
)

func TestInitConfig(t *testing.T) {
	fs := afero.NewMemMapFs()

	cfg := Config{
		Router:         "app",
		Language:       "ts",
		ComponentStyle: "const",
	}

	WriteConfig(fs, cfg)

	// Check file exists
	exists, err := afero.Exists(fs, ConfigFileName)
	if err != nil {
		t.Fatalf("Error checking if config file exists: %v", err)
	}
	if !exists {
		t.Fatalf("Expected config file '%s' to exist", ConfigFileName)
	}

	// Check contents
	data, err := afero.ReadFile(fs, ConfigFileName)
	if err != nil {
		t.Fatalf("Error reading config file: %v", err)
	}

	var readCfg Config
	if err := json.Unmarshal(data, &readCfg); err != nil {
		t.Fatalf("Error unmarshalling config JSON: %v", err)
	}

	if readCfg != cfg {
		t.Errorf("Config mismatch.\nExpected: %+v\nGot: %+v", cfg, readCfg)
	}
}
