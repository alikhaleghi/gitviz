package tui

import (
	"fmt"
	"os"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

// Model owns UI state for the initial scaffold.
type Model struct {
	width        int
	height       int
	path         string
	repoDetected bool
	view         string
	status       string
}

// NewModel builds a default model with lightweight runtime context.
func NewModel() Model {
	cwd, err := os.Getwd()
	if err != nil {
		cwd = "unknown"
	}

	_, statErr := os.Stat(".git")
	repoDetected := statErr == nil
	status := "Status: waiting for repository data"
	if !repoDetected {
		status = "Status: no repository detected"
	}

	return Model{
		path:         cwd,
		repoDetected: repoDetected,
		view:         "commits",
		status:       status,
	}
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
	if m.width == 0 {
		m.width = 100
	}
	if m.height == 0 {
		m.height = 30
	}

	frame := FrameStyle.Width(m.width - 2)

	contentWidth := frame.GetWidth() - frame.GetHorizontalPadding()
	if contentWidth < 40 {
		contentWidth = 40
	}

	header := m.renderHeader(contentWidth)
	main := m.renderMain(contentWidth)
	status := m.renderStatus(contentWidth)
	footer := m.renderFooter(contentWidth)

	body := lipgloss.JoinVertical(lipgloss.Left, header, main, status, footer)
	return frame.Render(body) + "\n"
}

func Run() error {
	p := tea.NewProgram(NewModel(), tea.WithAltScreen())
	_, err := p.Run()
	return err
}

func (m Model) renderHeader(width int) string {
	title := StyleAccent.Render("gitviz")
	repoState := "none"
	if m.repoDetected {
		repoState = "detected"
	}

	meta := fmt.Sprintf("Repo: %s  |  View: %s", repoState, m.view)
	path := fmt.Sprintf("Path: %s", m.path)
	line := strings.Repeat("─", width)

	return lipgloss.JoinVertical(
		lipgloss.Left,
		fmt.Sprintf("%s  %s", title, StyleMuted.Render(meta)),
		StyleMuted.Render(path),
		StyleMuted.Render(line),
	)
}

func (m Model) renderMain(width int) string {
	colGap := 1
	leftWidth := (width - colGap) / 2
	rightWidth := width - colGap - leftWidth

	panelTitle := StyleBold
	panel := StylePanel

	leftBody := "- (placeholder)"
	rightBody := "Select a commit to inspect."
	if !m.repoDetected {
		leftBody = "No Git repository detected\n\nThis directory is not initialized.\n\nTry: git init"
		rightBody = "Welcome to gitviz\n\nTo get started here:\n1. git init\n2. git add .\n3. git commit -m \"chore: initialize project scaffold\""
	}

	left := panel.Width(leftWidth).Render(
		lipgloss.JoinVertical(
			lipgloss.Left,
			panelTitle.Render("Commits"),
			strings.Repeat("─", Max(8, leftWidth-2)),
			leftBody,
		),
	)
	right := panel.Width(rightWidth).Render(
		lipgloss.JoinVertical(
			lipgloss.Left,
			panelTitle.Render("Details"),
			strings.Repeat("─", Max(8, rightWidth-2)),
			rightBody,
		),
	)

	line := StyleMuted.Render(strings.Repeat("─", width))
	panes := lipgloss.JoinHorizontal(lipgloss.Top, left, strings.Repeat(" ", colGap), right)
	return lipgloss.JoinVertical(lipgloss.Left, panes, line)
}

func (m Model) renderStatus(width int) string {
	line := StyleMuted.Render(strings.Repeat("─", width))
	return lipgloss.JoinVertical(lipgloss.Left, m.status, line)
}

func (m Model) renderFooter(width int) string {
	help := "q quit  r refresh  ↑↓ move  enter inspect  b branches"
	return StyleMuted.Render(lipgloss.NewStyle().Width(width).Render(help))
}
