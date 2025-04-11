package cmd

import (
	"bytes"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"text/template"

	"github.com/bllakcn/nextjs-routing-helper-cli/cmd/constants"
	"github.com/bllakcn/nextjs-routing-helper-cli/cmd/helpers"
	"github.com/spf13/afero"
	"github.com/spf13/cobra"
)

// addCmd represents the add command
var addCmd = &cobra.Command{
	Use:   "add [page-name] --flag",
	Short: "Adds a new page to your Next.js project.",
	Long: `Adds a new page based on the configuration.
Page name can include subdirectories (e.g., 'users/profile').`,
	Args: cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		pageNameInput := args[0]

		useClientFlag, _ := cmd.Flags().GetBool("use-client")

		// 1. Read Configuration
		config, err := loadConfig()
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error loading configuration:\n%v\n", err)
			fmt.Fprintln(os.Stderr, "Please run 'nextjs-routing-helper-cli init' first.")
			os.Exit(1)
		}

		// Check if the router type is valid
		if !config.Router.IsValid() {
			fmt.Fprintf(os.Stderr, "Invalid router type in configuration file '%s'.\n", constants.ConfigFileName)
			os.Exit(1)
		}

		// 2. Determine Path and Filename
		targetPath, pageComponentName, err := determinePathAndComponent(pageNameInput, config)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error determining path:\n%v\n", err)
			os.Exit(1)
		}

		// 3. Generate File Content
		content, err := generatePageContent(pageComponentName, config, useClientFlag)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error generating page content:\n%v\n", err)
			os.Exit(1)
		}

		// 4. Create Directories and File
		err = createPageFile(afero.NewOsFs(), targetPath, content)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error creating page file:\n%v\n", err)
			os.Exit(1)
		}

		fmt.Printf("Successfully created page:\n%s\n", targetPath)
	},
}

// loadConfig reads and parses the .nextjs_routing_helper.json file
func loadConfig() (*Config, error) {
	data, err := os.ReadFile(constants.ConfigFileName)
	if err != nil {
		return nil, fmt.Errorf("could not read config file '%s': %w", constants.ConfigFileName, err)
	}

	var config Config
	err = json.Unmarshal(data, &config)
	if err != nil {
		return nil, fmt.Errorf("could not parse config file '%s': %w", constants.ConfigFileName, err)
	}
	return &config, nil
}

// determinePathAndComponent calculates the final file path and component name
func determinePathAndComponent(pageNameInput string, config *Config) (filePath string, componentName string, err error) {
	parts := strings.Split(pageNameInput, "/")
	if len(parts) == 0 || parts[0] == "" {
		return "", "", fmt.Errorf("page name cannot be empty or just slashes")
	}
	// Ignore the last part if it's "index" for pages router
	if config.Router == constants.PagesRouter && strings.ToLower(parts[len(parts)-1]) == "index" {
		parts = parts[:len(parts)-1]
	}
	// Ignore the last part if it's "page" for app router
	if config.Router == constants.AppRouter && strings.ToLower(parts[len(parts)-1]) == "page" {
		parts = parts[:len(parts)-1]
	}
	fileNamePart := parts[len(parts)-1]
	if fileNamePart == "" {
		return "", "", fmt.Errorf("page name cannot end with a slash")
	}

	// Determine component name (e.g., "UserProfilePage", "AboutPage")
	componentName = helpers.ToPascalCase(fileNamePart)
	if config.PageComponentSuffix != "" {
		componentName += helpers.ToPascalCase(config.PageComponentSuffix)
	}

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
	if config.Router == constants.AppRouter {
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
		// For pages router, always use 'index.ext' as the page file
		pageFileName = "index" + fileExtension
		filePath = filepath.Join(basePath, pageNameInput, pageFileName)
	}

	// Clean the path (removes redundant slashes, resolves "..")
	filePath = filepath.Clean(filePath)

	return filePath, componentName, nil
}

// PageData holds the dynamic data for the page template
type PageData struct {
	ComponentName string
	Style         string
	UseClient     bool
}

// generatePageContent creates the basic component code
func generatePageContent(componentName string, config *Config, useClient bool) (string, error) {
	// Define the template
	const tpl = `
{{if .UseClient}}'use client';{{end}}
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
		UseClient:     config.Router == constants.AppRouter && useClient,
	}

	// Execute the template
	var output bytes.Buffer
	if err := tmpl.Execute(&output, data); err != nil {
		return "", err
	}

	return output.String(), nil
}

// createPageFile ensures directories exist and writes the file
func createPageFile(fs afero.Fs, targetPath string, content string) error {
	// Ensure directory exists
	dir := filepath.Dir(targetPath)
	if err := fs.MkdirAll(dir, 0755); err != nil { // rwxr-xr-x permissions
		return fmt.Errorf("could not create directory '%s': %w", dir, err)
	}

	// Write the file
	if err := afero.WriteFile(fs, targetPath, []byte(content), 0644); err != nil { // rw-r--r--
		return fmt.Errorf("could not write file '%s': %w", targetPath, err)
	}

	return nil
}

func init() {
	rootCmd.AddCommand(addCmd)
	addCmd.Flags().Bool("use-client", false, "Use 'use client' directive for the component (only for app router)")
}
