package cmd

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/bllakcn/nextjs-routing-helper-cli/cmd/constants"
	"github.com/spf13/afero"
	"github.com/spf13/cobra"
)

var viewCmd = &cobra.Command{
	Use:   "view",
	Short: "Visualizes the routes in your Next.js project.",
	Long:  `Scans the app/pages directory and prints out a tree of all available routes.`,
	Run: func(cmd *cobra.Command, args []string) {
		config, err := constants.LoadConfig()
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error loading configuration: %v\n", err)
			os.Exit(1)
		}

		startPath := "pages"
		if config.Router == constants.AppRouter {
			startPath = "app"
		}

		err = afero.Walk(AppFs, startPath, func(path string, info os.FileInfo, err error) error {
			if err != nil {
				return err
			}
			if info.IsDir() {
				return nil
			}
			if strings.HasSuffix(info.Name(), ".tsx") || strings.HasSuffix(info.Name(), ".jsx") || strings.HasSuffix(info.Name(), ".js") {
				relPath := strings.TrimPrefix(path, startPath)
				routePath := convertToRoute(relPath)
				fmt.Println(routePath)
			}
			return nil
		})

		if err != nil {
			fmt.Fprintf(os.Stderr, "Error walking the directory: %v\n", err)
			os.Exit(1)
		}
	},
}

// convertToRoute transforms a file path into a route-like format
func convertToRoute(filePath string) string {
	route := strings.TrimSuffix(filePath, filepath.Ext(filePath))
	route = strings.ReplaceAll(route, "\\", "/")
	route = strings.ReplaceAll(route, "/page", "")
	route = strings.ReplaceAll(route, "/index", "")
	if route == "" {
		return "/"
	}
	return route
}

func init() {
	rootCmd.AddCommand(viewCmd)
}
