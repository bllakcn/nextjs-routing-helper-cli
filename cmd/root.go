package cmd

import (
	"os"

	"github.com/spf13/afero"
	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Version: "v0.0.1",
	Use:     "nextjs-routing-helper",
	Example: `
	- nextjs-routing-helper init 
	- nextjs-routing-helper add [pageName]
	`,
	Short: "Nextjs Routing Helper CLI - a simple CLI to create pages in Nextjs",
	Long: `Nextjs Routing Helper CLI is a fast way to create pages in your Nextjs project. It creates necessary files based on your preferences.
	`,
}

// Create a filesystem instance. For production, use the OS filesystem.
var AppFs afero.Fs = afero.NewOsFs()

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
}
