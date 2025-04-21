package cmd

import (
	"fmt"
	"os"

	"github.com/bllakcn/nextjs-routing-helper-cli/cmd/constants"
	treeui "github.com/bllakcn/nextjs-routing-helper-cli/cmd/ui/tree"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/savannahostrowski/tree-bubble"

	"github.com/spf13/cobra"
)

var viewCmd = &cobra.Command{
	Use:   "view",
	Short: "Visualizes the routes in your Next.js project.",
	Long:  `Scans the app/pages directory and prints out a tree of all available routes.`,
	Run: func(cmd *cobra.Command, args []string) {
		// Load config
		config, err := constants.LoadConfig()
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error loading configuration: %v\n", err)
			os.Exit(1)
		}

		// Determine the router type
		startPath := "pages"
		if config.Router == constants.AppRouter {
			startPath = "app"
		}

		// Create the nodes
		nodeTree := treeui.BuildRouteTree(AppFs, startPath)
		//

		m := treeui.New([]tree.Node{nodeTree})
		p := tea.NewProgram(m, tea.WithAltScreen())

		if _, err := p.Run(); err != nil {
			fmt.Fprintf(os.Stderr, "Error running program: %v\n", err)
			os.Exit(1)
		}
	},
}

func init() {
	rootCmd.AddCommand(viewCmd)
}
