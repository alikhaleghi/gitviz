package tui

import (
	"strings"

	"github.com/charmbracelet/lipgloss"
)

func (m Model) renderDetails(width int) string {
	body := "Select a commit to inspect."
	if m.detailContent != "" {
		body = m.detailContent
	} else if !m.repoDetected {
		body = "Welcome to gitviz\n\nTo get started here:\n1. git init\n2. git add .\n3. git commit -m \"chore: initialize project scaffold\""
	}
	return StylePanel.Width(width).Render(
		lipgloss.JoinVertical(
			lipgloss.Left,
			StyleBold.Render("Details"),
			strings.Repeat("─", Max(8, width-2)),
			body,
		),
	)
}
