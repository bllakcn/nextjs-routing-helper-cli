package nextjs_routing_helper

import (
	"bufio"
	"encoding/json"
	"fmt"
	"os"
	"strings"
)

// Config holds the user preferences
type Config struct {
	Router         string `json:"router"`         // "app" or "pages"
	Language       string `json:"language"`       // "js" or "ts"
	ComponentStyle string `json:"componentStyle"` // "const" or "function"
	// Add more preferences as needed, e.g., namingConvention, useJsx
}

const configFileName = ".nextjs_routing_helper_cli.json" // Name of the config file

func Init(){
	reader := bufio.NewReader(os.Stdin)
	config := Config{}

	fmt.Println("Initializing Next.js Routing CLI configuration...")

		// Check if file exists
		if _, err := os.Stat(configFileName); err == nil {
			fmt.Printf("Configuration file '%s' already exists.\n", configFileName)
			fmt.Print("Overwrite? (y/N): ")
			overwriteChoice, _ := reader.ReadString('\n')
			if !strings.EqualFold(strings.TrimSpace(overwriteChoice), "y") {
				fmt.Println("Initialization cancelled.")
				return
			}
		}

	// --- Router Type ---
	fmt.Print("Use App Router or Pages Router? (app/pages): ")
	routerChoice, _ := reader.ReadString('\n')
	routerChoice = strings.ToLower(strings.TrimSpace(routerChoice))
	if routerChoice != "app" && routerChoice != "pages" {
		fmt.Println("Invalid choice. Defaulting to 'app'.")
		config.Router = "app"
	} else {
		config.Router = routerChoice
	}

	// --- Language ---
	fmt.Print("Use TypeScript or JavaScript? (ts/js): ")
	langChoice, _ := reader.ReadString('\n')
	langChoice = strings.ToLower(strings.TrimSpace(langChoice))
	if langChoice != "ts" && langChoice != "js" {
		fmt.Println("Invalid choice. Defaulting to 'ts'.")
		config.Language = "ts"
	} else {
		config.Language = langChoice
	}

	// --- Component Style ---
	fmt.Print("Prefer 'function' declarations or 'const' arrow functions? (function/const): ")
	styleChoice, _ := reader.ReadString('\n')
	styleChoice = strings.ToLower(strings.TrimSpace(styleChoice))
	if styleChoice != "const" && styleChoice != "function" {
		fmt.Println("Invalid choice. Defaulting to 'function'.")
		config.ComponentStyle = "function"
	} else {
		config.ComponentStyle = styleChoice
	}


	// Marshal config to JSON
	configData, err := json.MarshalIndent(config, "", "  ") // Pretty print JSON
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error marshalling config to JSON: %v\n", err)
		os.Exit(1)
	}

	// Write config file
	err = os.WriteFile(configFileName, configData, 0644) // rw-r--r-- permissions
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error writing config file '%s': %v\n", configFileName, err)
		os.Exit(1)
	}

	fmt.Printf("Configuration saved successfully to %s\n", configFileName)
	fmt.Printf("  Router: %s\n", config.Router)
	fmt.Printf("  Language: %s\n", config.Language)
	fmt.Printf("  Component Style: %s\n", config.ComponentStyle)

}