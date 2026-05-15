package tui

import (
	"strings"

	"github.com/charmbracelet/lipgloss"
)

func (m Model) renderCommitList(width int) string {
	if !m.repoDetected {
		return StylePanel.Width(width).Render(
			lipgloss.JoinVertical(
				lipgloss.Left,
				StyleBold.Render("Commits"),
				strings.Repeat("─", Max(8, width-2)),
				"No Git repository detected\n\nThis directory is not initialized.\n\nTry: git init",
			),
		)
	}

	titleStyle := StyleBold
	if m.focus == "commits" {
		titleStyle = StyleAccent
	}

	var items []string
	for i, commit := range m.commits {
		prefix := "  "
		style := StyleNormal
		if i == m.cursor {
			prefix = "▸ "
			style = StyleSelected
		}
		items = append(items, prefix+style.Render(commit))
	}
	body := strings.Join(items, "\n")

	return StylePanel.Width(width).Render(
		lipgloss.JoinVertical(
			lipgloss.Left,
			titleStyle.Render("Commits"),
			strings.Repeat("─", Max(8, width-2)),
			body,
		),
	)
}
