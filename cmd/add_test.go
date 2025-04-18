package cmd

import (
	"path/filepath"
	"testing"

	"github.com/bllakcn/nextjs-routing-helper-cli/cmd/constants"
	"github.com/stretchr/testify/assert"
)

func TestDeterminePathAndComponent(t *testing.T) {
	tests := []struct {
		configRouter              constants.RouterType
		configLanguage            constants.LanguageType
		configComponentStyle      constants.ComponentStyleType
		configPageComponentSuffix string
		configSrcFolder           bool
		inputPath                 string
		expectedTarget            string
		expectedName              string
	}{
		{
			configRouter:              "app",
			configLanguage:            "ts",
			configComponentStyle:      "function",
			configSrcFolder:           false,
			configPageComponentSuffix: "page",

			inputPath:      "dashboard",
			expectedTarget: filepath.Join("app", "dashboard", "page.tsx"),
			expectedName:   "DashboardPage",
		},
		{
			configRouter:              "pages",
			configLanguage:            "ts",
			configComponentStyle:      "const",
			configSrcFolder:           false,
			configPageComponentSuffix: "page",

			inputPath:      "auth/login",
			expectedTarget: filepath.Join("pages", "auth", "login", "index.tsx"),
			expectedName:   "LoginPage",
		},
		{
			configRouter:              "app",
			configLanguage:            "js",
			configComponentStyle:      "function",
			configSrcFolder:           true,
			configPageComponentSuffix: "component",

			inputPath:      "products/details",
			expectedTarget: filepath.Join("src", "app", "products", "details", "page.jsx"),
			expectedName:   "DetailsComponent",
		},
	}

	for _, tt := range tests {
		mockConfig := &constants.Config{
			Router:              tt.configRouter,
			Language:            tt.configLanguage,
			ComponentStyle:      tt.configComponentStyle,
			SrcFolder:           tt.configSrcFolder,
			PageComponentSuffix: tt.configPageComponentSuffix,
		}
		targetPath, componentName, err := determinePathAndComponent(tt.inputPath, mockConfig)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		assert.Equal(t, tt.expectedTarget, targetPath, "unexpected target path")
		assert.Equal(t, tt.expectedName, componentName, "unexpected component name")
	}
}
