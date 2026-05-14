package tui

import (
	"fmt"
	"os"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

// Model owns UI state for the initial scaffold.
type Model struct {
	width  int
	height int
	status string
}

// NewModel builds a default model with lightweight runtime context.
func NewModel() Model {
	cwd, err := os.Getwd()
	if err != nil {
		cwd = "unknown"
	}
	return Model{status: fmt.Sprintf("repo: %s", cwd)}
}

func (m Model) Init() tea.Cmd {
	return nil
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "q", "ctrl+c":
			return m, tea.Quit
		}
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
	}

	return m, nil
}

func (m Model) View() string {
	title := lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("12")).Render("gitviz")
	header := fmt.Sprintf("%s  |  %s", title, m.status)

	leftPane := lipgloss.NewStyle().Border(lipgloss.RoundedBorder()).Padding(0, 1).Render("Commits\n\n- (placeholder)")
	rightPane := lipgloss.NewStyle().Border(lipgloss.RoundedBorder()).Padding(0, 1).Render("Details\n\nSelect a commit to inspect.")
	footer := lipgloss.NewStyle().Foreground(lipgloss.Color("8")).Render("q: quit | arrows: navigate (next)")

	panes := lipgloss.JoinHorizontal(lipgloss.Top, leftPane, rightPane)
	return lipgloss.JoinVertical(lipgloss.Left, header, "", panes, "", footer) + "\n"
}

func Run() error {
	p := tea.NewProgram(NewModel(), tea.WithAltScreen())
	_, err := p.Run()
	return err
}
