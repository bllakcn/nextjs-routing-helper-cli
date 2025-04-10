package cmd

import (
	"bytes"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"text/template"

	"github.com/bllakcn/nextjs-routing-helper-cli/pkg"
	"github.com/spf13/cobra"
)

// addCmd represents the add command
var addCmd = &cobra.Command{
	Use:   "add <page-name>",
	Short: "Adds a new page to your Next.js project.",
	Long: `Adds a new page based on the configuration.
Page name can include subdirectories (e.g., 'users/profile').`,
	Args: cobra.ExactArgs(1), // Ensures exactly one argument (page-name) is provided
	Run: func(cmd *cobra.Command, args []string) {
		pageNameInput := args[0]

		// 1. Read Configuration
		config, err := loadConfig()
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error loading configuration:\n%v", err)
			fmt.Fprintln(os.Stderr, "Please run 'nextjs-routing-helper-cli init' first.")
			os.Exit(1)
		}

		// 2. Determine Path and Filename
		targetPath, pageComponentName, err := determinePathAndComponent(pageNameInput, config)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error determining path:\n%v\n", err)
			os.Exit(1)
		}

		// 3. Generate File Content
		content, err := generatePageContent(pageComponentName, config)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error generating page content:\n%v\n", err)
			os.Exit(1)
		}

		// 4. Create Directories and File
		err = createPageFile(targetPath, content)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error creating page file:\n%v\n", err)
			os.Exit(1)
		}

		fmt.Printf("Successfully created page:\n%s\n", targetPath)
	},
}

// loadConfig reads and parses the .nextjs_routing_helper.json file
func loadConfig() (*Config, error) {
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

// determinePathAndComponent calculates the final file path and component name
func determinePathAndComponent(pageNameInput string, config *Config) (filePath string, componentName string, err error) {
	parts := strings.Split(pageNameInput, "/")
	if len(parts) == 0 || parts[0] == "" {
		return "", "", fmt.Errorf("page name cannot be empty or just slashes")
	}
	fileNamePart := parts[len(parts)-1]
	if fileNamePart == "" {
		return "", "", fmt.Errorf("page name cannot end with a slash")
	}

	// Determine component name (e.g., "UserProfilePage", "AboutPage")
	componentName = pkg.ToPascalCase(fileNamePart) + "Page" // TODO: to be determined with config

	var basePath string
	var fileExtension string
	var pageFileName string

	// Determine file extension
	if config.Language == "ts" {
		fileExtension = ".tsx"
	} else {
		fileExtension = ".jsx"
	}

	// Determine base path and the actual filename used in the path structure
	if config.Router == "app" {
		if config.SrcFolder {
			basePath = filepath.Join("src", "app")
		} else {
			basePath = "app"
		}
		// App router always uses 'page.ext' in its leaf directory
		pageFileName = "page" + fileExtension
		// Construct path using the folder structure from input and the required filename
		filePath = filepath.Join(basePath, pageNameInput, pageFileName)
	} else { // pages router
		if config.SrcFolder {
			basePath = filepath.Join("src", "pages")
		} else {
			basePath = "pages"
		}
		// Handle index pages specifically
		if strings.ToLower(fileNamePart) == "index" {
			pageFileName = "index" + fileExtension
			if len(parts) > 1 { // Index inside a subdirectory e.g. users/index
				// Path is pages/users/index.tsx
				filePath = filepath.Join(basePath, filepath.Join(parts[:len(parts)-1]...), pageFileName)
			} else { // Top-level index
				// Path is pages/index.tsx
				filePath = filepath.Join(basePath, pageFileName)
			}
		} else { // Regular page file, e.g., "about" or "users/profile"
			pageFileName = fileNamePart + fileExtension
			// Path uses the input path structure directly, ending with the calculated filename
			// For "users/profile": pages/users/profile.tsx -> Join(basePath, "users", "profile.tsx")
			// For "about": pages/about.tsx -> Join(basePath, "about.tsx")
			parentDirs := parts[:len(parts)-1]
			filePath = filepath.Join(basePath, filepath.Join(parentDirs...), pageFileName)
		}
	}

	// Clean the path (removes redundant slashes, resolves "..")
	filePath = filepath.Clean(filePath)

	return filePath, componentName, nil
}

// PageData holds the dynamic data for the page template
type PageData struct {
	ComponentName string
	Style         string
	// UseClient     bool
}

// generatePageContent creates the basic component code
// THIS IS A VERY BASIC EXAMPLE - Use text/template for real world use
func generatePageContent(componentName string, config *Config) (string, error) {
	// Define the template
	// {{if .UseClient}}'use client';{{end}} //TODO ↴
	const tpl = `
{{if eq .Style "const"}}const {{.ComponentName}} = () => {
  return (
    <div>
      <h1>{{.ComponentName}}</h1>
      {/* Add your content here */}
    </div>
  );
};

export default {{.ComponentName}};
{{else}}export default function {{.ComponentName}}() {
  return (
    <div>
      <h1>{{.ComponentName}}</h1>
      {/* Add your content here */}
    </div>
  );
}
{{end}}`

	// Parse the template
	tmpl, err := template.New("page").Parse(tpl)
	if err != nil {
		return "", err
	}

	// Prepare the data
	data := PageData{
		ComponentName: componentName,
		Style:         config.ComponentStyle,
		// UseClient:     config.Router == "app" &&  //TODO: This should be a flag
	}

	// Execute the template
	var output bytes.Buffer
	if err := tmpl.Execute(&output, data); err != nil {
		return "", err
	}

	return output.String(), nil
}

// createPageFile ensures directories exist and writes the file
func createPageFile(targetPath string, content string) error {
	// Ensure directory exists
	dir := filepath.Dir(targetPath)
	if err := os.MkdirAll(dir, 0755); err != nil { // rwxr-xr-x permissions
		return fmt.Errorf("could not create directory '%s': %w", dir, err)
	}

	// Write the file
	if err := os.WriteFile(targetPath, []byte(content), 0644); err != nil { // rw-r--r--
		return fmt.Errorf("could not write file '%s': %w", targetPath, err)
	}

	return nil
}

func init() {
	rootCmd.AddCommand(addCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// addCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// addCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
