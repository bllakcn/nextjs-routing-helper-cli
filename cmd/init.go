/*
Copyright © 2025 AHMET BILAL AKCAN <bllakcn35@gmail.com>
*/
package cmd

import (
	"fmt"

	nextjs_routing "github.com/bllakcn/nextjs-routing-helper-cli/pkg"
	"github.com/spf13/cobra"
)

const configFileName = ".nextjs_routing_helper_cli.json" // Name of the config file


// initCmd represents the init command
var initCmd = &cobra.Command{
	Use:   "init",
	Short: "Initialize your preferences",
	Long: fmt.Sprintf(`Creates a %s file in the current directory
to store your preferences for generating Next.js pages and components.

This file contains settings like:
- Router type (app/pages)
- Language (ts/js)
- Component style (const/function)

If the file already exists, you will be prompted to overwrite it.`, configFileName),
	Run: func(cmd *cobra.Command, args []string) {
		nextjs_routing.Init()
	},
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
