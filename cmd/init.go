package cmd

import (
	"bufio"
	"encoding/json"
	"fmt"
	"os"
	"strings"

	"github.com/spf13/afero"
	"github.com/spf13/cobra"
)

// Create a filesystem instance. For production, use the OS filesystem.
var AppFs afero.Fs = afero.NewOsFs()

const ConfigFileName = ".nextjs_routing_helper.json"

// Config holds the user preferences
type Config struct {
	Router         string `json:"router"`         // "app" or "pages"
	Language       string `json:"language"`       // "js" or "ts"
	ComponentStyle string `json:"componentStyle"` // "function" or "const"
	SrcFolder      bool   `json:"srcFolder"`
}

var initCmd = &cobra.Command{
	Use:   "init",
	Short: "Initialize your preferences",
	Long: fmt.Sprintf(`Creates a %s file in the current directory
to store your preferences for generating Next.js pages and components.

This file contains settings like:
- Router type (app/pages)
- Language (ts/js)
- Component style (const/function)
- Src folder (yes/no)

If the file already exists, you will be prompted to overwrite it.`, ConfigFileName),
	Run: func(cmd *cobra.Command, args []string) {
		reader := bufio.NewReader(os.Stdin)
		config := Config{}

		fmt.Println("Initializing Next.js Routing CLI configuration...")

		// Check if file exists
		if exists, _ := afero.Exists(AppFs, ConfigFileName); exists {
			fmt.Printf("Configuration file '%s' already exists.\n", ConfigFileName)
			fmt.Print("Overwrite? (y/N): ")
			overwriteChoice, _ := reader.ReadString('\n')
			if !strings.EqualFold(strings.TrimSpace(overwriteChoice), "y") {
				fmt.Println("Initialization cancelled.")
				return
			}
		}

		// --- Router Type ---
		fmt.Print("Use App Router or Pages Router? (App/pages): ")
		routerChoice, _ := reader.ReadString('\n')
		routerChoice = strings.ToLower(strings.TrimSpace(routerChoice))
		if routerChoice != "app" && routerChoice != "pages" {
			fmt.Println("Invalid choice. Defaulting to 'app'.")
			config.Router = "app"
		} else {
			config.Router = routerChoice
		}

		// --- Src Folder ---
		fmt.Print("Use `src/` folder? (No/yes): ")
		srcChoice, _ := reader.ReadString('\n')
		srcChoice = strings.ToLower(strings.TrimSpace(srcChoice))
		if srcChoice != "yes" && srcChoice != "no" {
			fmt.Println("Invalid choice. Defaulting to no.")
			config.SrcFolder = false
		} else {
			if srcChoice == "yes" {
				config.SrcFolder = true
			} else {
				config.SrcFolder = false
			}
		}

		// --- Language ---
		fmt.Print("Use TypeScript or JavaScript? (Ts/js): ")
		langChoice, _ := reader.ReadString('\n')
		langChoice = strings.ToLower(strings.TrimSpace(langChoice))
		if langChoice != "ts" && langChoice != "js" {
			fmt.Println("Invalid choice. Defaulting to 'ts'.")
			config.Language = "ts"
		} else {
			config.Language = langChoice
		}

		// --- Component Style ---
		fmt.Print("Prefer 'function' declarations or 'const' arrow functions? (Function/const): ")
		styleChoice, _ := reader.ReadString('\n')
		styleChoice = strings.ToLower(strings.TrimSpace(styleChoice))
		if styleChoice != "const" && styleChoice != "function" {
			fmt.Println("Invalid choice. Defaulting to 'function'.")
			config.ComponentStyle = "function"
		} else {
			config.ComponentStyle = styleChoice
		}

		WriteConfig(AppFs, config)

		fmt.Printf("Configuration saved successfully to %s\n", ConfigFileName)
		fmt.Printf("  Router: %s\n", config.Router)
		fmt.Printf("  Language: %s\n", config.Language)
		fmt.Printf("  Component Style: %s\n", config.ComponentStyle)

	},
}

// WriteConfig writes the config to the given filesystem.
func WriteConfig(fs afero.Fs, config Config) error {
	// Marshal config to JSON
	configData, err := json.MarshalIndent(config, "", "  ") // Pretty print JSON
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error marshalling config to JSON: %v\n", err)
		os.Exit(1)
	}

	// Write config file
	err = afero.WriteFile(fs, ConfigFileName, configData, 0644) // rw-r--r-- permissions
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error writing config file '%s': %v\n", ConfigFileName, err)
		os.Exit(1)
	}
	return err
}

func init() {
	rootCmd.AddCommand(initCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// initCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// initCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
