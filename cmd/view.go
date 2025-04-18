package cmd

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/bllakcn/nextjs-routing-helper-cli/cmd/constants"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/charmbracelet/x/term"
	tree "github.com/savannahostrowski/tree-bubble"
	"github.com/spf13/afero"
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

		startPath := "pages"
		if config.Router == constants.AppRouter {
			startPath = "app"
		}

		rootNode := buildRouteTree(startPath)
		w, h, _ := term.GetSize(os.Stdout.Fd())
		p := tea.NewProgram(initialModel([]tree.Node{rootNode}, w, h))

		if _, err := p.Run(); err != nil {
			fmt.Printf("Alas, there's been an error: %v", err)
			os.Exit(1)
		}
	},
}

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

var (
	styleDoc = lipgloss.NewStyle().
		PaddingTop(2).
		PaddingLeft(4)
)

var accentColor = lipgloss.Color("#f3bd72")

type model struct {
	tree tree.Model
}

func (m model) Init() tea.Cmd {
	return nil
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "q", "ctrl+c":
			return m, tea.Quit
		}
	}
	var cmd tea.Cmd
	m.tree, cmd = m.tree.Update(msg)
	return m, cmd
}

func (m model) View() string {
	return styleDoc.Render(m.tree.View())
}

func initialModel(nodes []tree.Node, w, h int) model {
	m := tree.New(nodes, w, h)
	m.Styles.Selected = lipgloss.NewStyle().Foreground(accentColor)
	m.Styles.Shapes = lipgloss.NewStyle().Foreground(accentColor)

	return model{tree: m}
}

func buildRouteTree(startPath string) tree.Node {
	nodeMap := make(map[string]*tree.Node)

	_ = afero.Walk(AppFs, startPath, func(path string, info os.FileInfo, err error) error {
		if err != nil || info.IsDir() {
			return nil
		}
		if strings.HasSuffix(info.Name(), ".tsx") || strings.HasSuffix(info.Name(), ".jsx") || strings.HasSuffix(info.Name(), ".js") {
			relPath := strings.TrimPrefix(path, startPath)
			routePath := convertToRoute(relPath)
			parts := strings.Split(strings.Trim(routePath, "/"), "/")

			var fullPath string
			for i, part := range parts {
				if i == 0 {
					fullPath = part
				} else {
					fullPath = filepath.Join(fullPath, part)
				}
				if _, exists := nodeMap[fullPath]; !exists {
					relFile := strings.TrimPrefix(path, string(filepath.Separator))
					desc := relFile
					nodeMap[fullPath] = &tree.Node{
						Value:    part,
						Desc:     desc,
						Children: []tree.Node{},
					}
				}

				if i > 0 {
					parentPath := filepath.Join(parts[:i]...)
					parent := nodeMap[parentPath]
					child := nodeMap[fullPath]
					alreadyChild := false
					for _, c := range parent.Children {
						if c.Value == child.Value {
							alreadyChild = true
							break
						}
					}
					if !alreadyChild {
						parent.Children = append(parent.Children, *child)
					}
				}
			}
		}
		return nil
	})

	rootValue := strings.TrimPrefix(startPath, string(filepath.Separator))
	root, ok := nodeMap[rootValue]
	if !ok {
		root = &tree.Node{Value: rootValue}
	} else {
		root.Value = rootValue
	}

	// Collect all root-level children
	var rootChildren []tree.Node
	for key, node := range nodeMap {
		if !strings.Contains(key, string(filepath.Separator)) {
			rootChildren = append(rootChildren, *node)
		}
	}
	root.Children = rootChildren
	return *root
}

func init() {
	rootCmd.AddCommand(viewCmd)
}
