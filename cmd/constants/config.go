package constants

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/spf13/afero"
)

const ConfigFileName = ".nextjs_routing_helper.json"

type Config struct {
	Router              RouterType         `json:"router"`
	Language            LanguageType       `json:"language"`
	ComponentStyle      ComponentStyleType `json:"componentStyle"`
	SrcFolder           bool               `json:"srcFolder"`
	PageComponentSuffix string             `json:"pageComponentSuffix"`
}

// loadConfig reads and parses the config file
func LoadConfig() (*Config, error) {
	data, err := os.ReadFile(ConfigFileName)
	if err != nil {
		return nil, fmt.Errorf("could not read config file '%s': %w", ConfigFileName, err)
	}

	var config Config
	err = json.Unmarshal(data, &config)
	if err != nil {
		return nil, fmt.Errorf("could not parse config file '%s': %w", ConfigFileName, err)
	}
	return &config, nil
}

// WriteConfig writes the config to the given filesystem.
func WriteConfig(fs afero.Fs, config Config) error {
	// Marshal config to JSON
	configData, err := json.MarshalIndent(config, "", "  ") // Pretty print JSON
	if err != nil {
		return fmt.Errorf("error marshalling config to JSON: %w", err)
	}

	// Write config file
	err = afero.WriteFile(fs, ConfigFileName, configData, 0644) // rw-r--r-- permissions
	if err != nil {
		return fmt.Errorf("error writing config file '%s': %w", ConfigFileName, err)
	}
	return nil
}
