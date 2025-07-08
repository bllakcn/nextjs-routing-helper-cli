package cmd

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strings"
	"time"

	"github.com/bllakcn/nextjs-routing-helper-cli/cmd/constants"
	"github.com/bllakcn/nextjs-routing-helper-cli/cmd/helpers"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/common-nighthawk/go-figure"
	"github.com/spf13/afero"
	"github.com/spf13/cobra"
)

type Styles struct {
	InputField lipgloss.Style
}

func DefaultStyles() *Styles {
	s := new(Styles)
	s.InputField = lipgloss.NewStyle().
		Width(80).
		BorderForeground(constants.AccentColor).
		BorderStyle(lipgloss.NormalBorder())
	return s
}

type Question struct {
	question    string
	placeholder string
	answer      string
}

func NewQuestion(question string, placeholder string) Question {
	return Question{question: question, placeholder: placeholder}
}

type model struct {
	questions   []Question
	index       int
	answerField textinput.Model
	width       int
	height      int
	styles      *Styles
	done        bool
}

func New(questions []Question) *model {
	styles := DefaultStyles()
	answerField := textinput.New()
	answerField.Placeholder = questions[0].placeholder
	answerField.Width = 80
	answerField.Focus()
	return &model{questions: questions, answerField: answerField, styles: styles}
}

func (m *model) Next() {
	if m.index < len(m.questions)-1 {
		m.index++
		m.answerField.Placeholder = m.questions[m.index].placeholder
		m.answerField.Focus()
	}
}

func (m model) Init() tea.Cmd {
	return nil
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd
	current := &m.questions[m.index]
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.height = msg.Height
		m.width = msg.Width
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c":
			return m, tea.Quit
		case "enter":
			current.answer = m.answerField.Value()
			if m.index == len(m.questions)-1 {
				m.done = true
				return m, tea.Tick(time.Second*2, func(t time.Time) tea.Msg {
					return tea.Quit()
				})
			}
			m.Next()
			m.answerField.SetValue("")
			return m, nil
		}
	}
	m.answerField, cmd = m.answerField.Update(msg)
	return m, cmd
}

func (m model) View() string {
	if m.done {
		return lipgloss.Place(
			m.width,
			m.height,
			lipgloss.Center,
			lipgloss.Center,
			lipgloss.NewStyle().
				Foreground(constants.AccentColor).
				Bold(true).
				Render(string(fmt.Append(
					[]byte(figure.NewFigure("Next.js Routing Helper", "rectangles", true).String()),
					fmt.Sprintf("Saving config to %s\n", constants.ConfigFileName)),
				)),
		)
	}
	if m.width == 0 {
		return "loading..."
	}
	return lipgloss.Place(
		m.width,
		m.height,
		lipgloss.Center,
		lipgloss.Center,
		lipgloss.JoinVertical(lipgloss.Center, m.questions[m.index].question, m.styles.InputField.Render(m.answerField.View())),
	)
}

var initCmd = &cobra.Command{
	Use:   "init",
	Short: "Initialize your preferences",
	Long: fmt.Sprintf(`Creates a %s file in the current directory
to store your preferences for generating Next.js pages and components.

This file contains settings like:
- Router type (app/pages)
- Src folder (yes/no)
- Language (ts/js)
- Component style (const/function)

If the file already exists, you will be prompted to overwrite it.`, constants.ConfigFileName),
	Run: func(cmd *cobra.Command, args []string) {

		reader := bufio.NewReader(os.Stdin)

		// Check if file exists
		if exists, _ := afero.Exists(AppFs, constants.ConfigFileName); exists {
			fmt.Printf("Configuration file '%s' already exists.\n", constants.ConfigFileName)
			fmt.Print("Overwrite? (y/N): ")
			overwriteChoice, _ := reader.ReadString('\n')
			if !strings.EqualFold(strings.TrimSpace(overwriteChoice), "y") {
				fmt.Println("Initialization cancelled.")
				return
			}
		}

		questions := []Question{
			NewQuestion(
				"Use App Router or Pages Router?",
				fmt.Sprintf("(%s/%s)", helpers.ToPascalCase(string(constants.AppRouter)), constants.PagesRouter),
			),
			NewQuestion(
				"Does your project use a 'src' directory?",
				"(y/N)",
			),
			NewQuestion(
				"Use TypeScript or JavaScript?",
				"(Ts/js)",
			),
			NewQuestion(
				"Prefer 'function' declarations or 'const' arrow functions?",
				"(Function/const)",
			),
			NewQuestion(
				"Use 'Page' suffix for page components?",
				"(Y/n)",
			),
		}

		m := New(questions)

		f, err := tea.LogToFile("debug.log", "debug")
		if err != nil {
			log.Fatalf("err: %s", err)
		}
		defer f.Close()

		p := tea.NewProgram(m, tea.WithAltScreen())
		if _, err := p.Run(); err != nil {
			log.Fatal(err)
		}

		// Map answers to config
		config := constants.Config{}
		for i, q := range m.questions {
			ans := strings.ToLower(strings.TrimSpace(q.answer))
			switch i {
			case 0:
				if ans == string(constants.AppRouter) {
					config.Router = constants.AppRouter
				} else if ans == string(constants.PagesRouter) {
					config.Router = constants.PagesRouter
				} else {
					config.Router = constants.AppRouter
				}
			case 1:
				config.SrcFolder = strings.EqualFold(ans, "y")
			case 2:
				if ans == "ts" {
					config.Language = "ts"
				} else if ans == "js" {
					config.Language = "js"
				} else {
					config.Language = "ts"
				}
			case 3:
				if ans == "const" {
					config.ComponentStyle = "const"
				} else if ans == "function" {
					config.ComponentStyle = "function"
				} else {
					config.ComponentStyle = "function"
				}
			case 4:
				if ans == "n" {
					config.PageComponentSuffix = ""
				} else {
					config.PageComponentSuffix = "page"
				}
			}
		}

		// Write config to file
		if err := constants.WriteConfig(AppFs, config); err != nil {
			fmt.Fprintf(os.Stderr, "Failed to save configuration: %v\n", err)
			os.Exit(1)
		}
	},
}

func init() {
	rootCmd.AddCommand(initCmd)
}
