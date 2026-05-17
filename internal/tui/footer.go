package tui

import (
	"strings"

	"github.com/charmbracelet/lipgloss"
)

func (m Model) renderStatus(width int) string {
	line := StyleMuted.Render(strings.Repeat("─", width))
	return lipgloss.JoinVertical(lipgloss.Left, m.status, line)
}

func (m Model) renderFooter(width int) string {
	help := "q quit  r refresh  ↑↓/w s move  y copy hash  enter/e inspect  b branches"
	if m.showModal {
		help = "↑↓ nav  esc close"
	} else if m.detailFocus == "diff" {
		help = "q quit  j/k scroll diff  Tab exit diff"
	} else if m.focus == "details" {
		help = "q quit  Tab focus diff/ list  esc back"
	}
	return StyleMuted.Render(lipgloss.NewStyle().Width(width).Render(help))
}
