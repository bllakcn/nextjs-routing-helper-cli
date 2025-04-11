package cmd

import (
	"encoding/json"
	"testing"

	"github.com/bllakcn/nextjs-routing-helper-cli/cmd/constants"
	"github.com/spf13/afero"
	"github.com/stretchr/testify/assert"
)

func TestInitConfig(t *testing.T) {
	fs := afero.NewMemMapFs()

	cfg := Config{
		Router:         "app",
		Language:       "ts",
		ComponentStyle: "const",
		SrcFolder:      true,
	}

	WriteConfig(fs, cfg)

	// Check file exists
	exists, err := afero.Exists(fs, constants.ConfigFileName)
	assert.NoError(t, err, "checking if config file exists")
	assert.True(t, exists, "Expected config file '%s' to exist", constants.ConfigFileName)

	// Check contents
	data, err := afero.ReadFile(fs, constants.ConfigFileName)
	assert.NoError(t, err, "reading config file")

	var readCfg Config
	assert.NoError(t, json.Unmarshal(data, &readCfg), "unmarshalling config JSON")

	assert.Equal(t, cfg, readCfg, "Config mismatch")
}
