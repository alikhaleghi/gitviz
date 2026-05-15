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
	help := "q quit  r refresh  ↑↓ move  enter inspect  b branches"
	if m.focus == "details" {
		help = "q quit  Tab focus list  esc back"
	}
	return StyleMuted.Render(lipgloss.NewStyle().Width(width).Render(help))
}
