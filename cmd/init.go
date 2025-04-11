package cmd

import (
	"bufio"
	"encoding/json"
	"fmt"
	"os"
	"strings"

	"github.com/bllakcn/nextjs-routing-helper-cli/cmd/constants"
	"github.com/bllakcn/nextjs-routing-helper-cli/cmd/helpers"
	"github.com/spf13/afero"
	"github.com/spf13/cobra"
)

// Create a filesystem instance. For production, use the OS filesystem.
var AppFs afero.Fs = afero.NewOsFs()

// Config holds the user preferences
type Config struct {
	Router              constants.RouterType `json:"router"`
	Language            string               `json:"language"`
	ComponentStyle      string               `json:"componentStyle"`
	SrcFolder           bool                 `json:"srcFolder"`
	PageComponentSuffix string               `json:"pageComponentSuffix"`
}

var initCmd = &cobra.Command{
	Use:   "init",
	Short: "Initialize your preferences",
	Long: fmt.Sprintf(`Creates a %s file in the current directory
to store your preferences for generating Next.js pages and components.

This file contains settings like:
- Router type (app/pages)
- Src folder (yes/no)
- Language (ts/js)
- Component style (const/function)

If the file already exists, you will be prompted to overwrite it.`, constants.ConfigFileName),
	Run: func(cmd *cobra.Command, args []string) {
		reader := bufio.NewReader(os.Stdin)
		config := Config{}

		fmt.Println("Initializing Next.js Routing CLI configuration...")

		// Check if file exists
		if exists, _ := afero.Exists(AppFs, constants.ConfigFileName); exists {
			fmt.Printf("Configuration file '%s' already exists.\n", constants.ConfigFileName)
			fmt.Print("Overwrite? (y/N): ")
			overwriteChoice, _ := reader.ReadString('\n')
			if !strings.EqualFold(strings.TrimSpace(overwriteChoice), "y") {
				fmt.Println("Initialization cancelled.")
				return
			}
		}

		// --- Router Type ---
		fmt.Printf("Use App Router or Pages Router? (%s/%s): ", helpers.ToPascalCase(string(constants.AppRouter)), constants.PagesRouter)
		routerChoiceStr, _ := reader.ReadString('\n') // Read user input as string
		routerChoiceStr = strings.ToLower(strings.TrimSpace(routerChoiceStr))
		// Validate the input string and assign the corresponding RouterType constant
		switch routerChoiceStr {
		case string(constants.AppRouter): // Compare input string to the string value of AppRouter ("app")
			config.Router = constants.AppRouter
		case string(constants.PagesRouter): // Compare input string to the string value of PagesRouter ("pages")
			config.Router = constants.PagesRouter
		default:
			// Handle invalid input - assign the default constant
			fmt.Printf("Invalid choice '%s'. Defaulting to '%s'.\n", routerChoiceStr, constants.AppRouter)
			config.Router = constants.AppRouter
		}

		// --- Src Folder ---
		fmt.Print("Does your project use a 'src' directory? (y/N): ")
		srcChoice, _ := reader.ReadString('\n')
		if strings.EqualFold(strings.ToLower(strings.TrimSpace(srcChoice)), "y") {
			config.SrcFolder = true
		} else {
			config.SrcFolder = false
		}

		// --- Language ---
		fmt.Print("Use TypeScript or JavaScript? (Ts/js): ")
		langChoice, _ := reader.ReadString('\n')
		langChoice = strings.ToLower(strings.TrimSpace(langChoice))
		switch langChoice {
		case "ts":
			config.Language = "ts"
		case "js":
			config.Language = "js"
		default:
			fmt.Printf("Invalid choice '%s'. Defaulting to 'ts'.\n", langChoice)
			config.Language = "ts"
		}

		// --- Component Style ---
		fmt.Print("Prefer 'function' declarations or 'const' arrow functions? (Function/const): ")
		styleChoice, _ := reader.ReadString('\n')
		styleChoice = strings.ToLower(strings.TrimSpace(styleChoice))
		switch styleChoice {
		case "const":
			config.ComponentStyle = "const"
		case "function":
			config.ComponentStyle = "function"
		default:
			fmt.Printf("Invalid choice '%s'. Defaulting to 'function'.\n", styleChoice)
			config.ComponentStyle = "function"
		}

		// --- Page Suffix ---
		fmt.Print("Use 'Page' suffix for page components? (Y/n): ")
		pageSuffixChoice, _ := reader.ReadString('\n')
		if strings.EqualFold(strings.ToLower(strings.TrimSpace(pageSuffixChoice)), "n") {
			config.PageComponentSuffix = ""
		} else {
			config.PageComponentSuffix = "page"
		}

		// --- Write Config ---
		if err := WriteConfig(AppFs, config); err != nil {
			fmt.Fprintf(os.Stderr, "Failed to save configuration: %v\n", err)
			os.Exit(1)

		}
		fmt.Printf("Configuration saved successfully to %s\n", constants.ConfigFileName)
		fmt.Printf("  Router: %s\n", config.Router)
		fmt.Printf("  Src Folder: %t\n", config.SrcFolder)
		fmt.Printf("  Language: %s\n", config.Language)
		fmt.Printf("  Component Style: %s\n", config.ComponentStyle)
		fmt.Printf("  Page Component Suffix: %s\n", config.PageComponentSuffix)

	},
}

// WriteConfig writes the config to the given filesystem.
func WriteConfig(fs afero.Fs, config Config) error {
	// Marshal config to JSON
	configData, err := json.MarshalIndent(config, "", "  ") // Pretty print JSON
	if err != nil {
		return fmt.Errorf("error marshalling config to JSON: %w", err)
	}

	// Write config file
	err = afero.WriteFile(fs, constants.ConfigFileName, configData, 0644) // rw-r--r-- permissions
	if err != nil {
		return fmt.Errorf("error writing config file '%s': %w", constants.ConfigFileName, err)
	}
	return nil
}

func init() {
	rootCmd.AddCommand(initCmd)
}
