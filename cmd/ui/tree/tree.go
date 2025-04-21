package treeui

import (
	"os"
	"path/filepath"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/charmbracelet/x/term"
	tree "github.com/savannahostrowski/tree-bubble"
	"github.com/spf13/afero"
)

var accentColor = lipgloss.Color("#f3bd72")

var (
	styleDoc = lipgloss.NewStyle().
		PaddingTop(2).
		PaddingLeft(2)
)

type Model struct {
	tree   tree.Model
	width  int
	height int
}

func (m Model) Init() tea.Cmd {
	return nil
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
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

func (m Model) View() string {
	if m.width == 0 {
		return "loading..."
	}
	return styleDoc.Render(m.tree.View())
}

func New(nodes []tree.Node) Model {
	w, h, _ := term.GetSize(os.Stdout.Fd())
	m := tree.New(nodes, w, h)

	m.Styles.Selected = lipgloss.NewStyle().Foreground(accentColor)
	m.Styles.Shapes = lipgloss.NewStyle().Foreground(accentColor)

	return Model{tree: m}
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

// Building a hierarchical tree representation of the route structure starting with the `startPath`.
func BuildRouteTree(fs afero.Fs, startPath string) tree.Node {
	nodeMap := make(map[string]*tree.Node)

	_ = afero.Walk(fs, startPath, func(path string, info os.FileInfo, err error) error {
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
					nodeMap[fullPath] = &tree.Node{
						Value:    part,
						Desc:     relFile,
						Children: []tree.Node{},
					}
				}

				// Link to the parent
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
