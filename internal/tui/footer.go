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
	help := "q quit  r refresh  ↑↓/w s move  Space compare  y copy hash  enter/e inspect  b branches"
	if m.showModal {
		help = "↑↓ nav  esc close"
	} else if m.compareMode {
		help = "q quit  j/k scroll compare  Esc clear  Space unselect"
	} else if len(m.compareSelected) == 1 {
		help = "q quit  Space select 2nd  ↑↓/w s move  enter/e inspect  b branches"
	} else if m.blameView {
		help = "q quit  j/k scroll blame  Esc back  B blame"
	} else if m.detailFocus == "diff" {
		help = "q quit  j/k scroll diff  Tab exit diff  B blame file"
	} else if m.focus == "details" {
		help = "q quit  Tab focus diff/ list  esc back  B blame file"
	}
	return StyleMuted.Render(lipgloss.NewStyle().Width(width).Render(help))
}
