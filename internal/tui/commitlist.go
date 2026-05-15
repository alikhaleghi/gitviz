package tui

import (
	"strings"

	"github.com/charmbracelet/lipgloss"
)

func (m Model) renderCommitList(width int) string {
	body := "- (placeholder)"
	if !m.repoDetected {
		body = "No Git repository detected\n\nThis directory is not initialized.\n\nTry: git init"
	} else {
		body = strings.Join(m.commits, "\n")
	}
	return StylePanel.Width(width).Render(
		lipgloss.JoinVertical(
			lipgloss.Left,
			StyleBold.Render("Commits"),
			strings.Repeat("─", Max(8, width-2)),
			body,
		),
	)
}
