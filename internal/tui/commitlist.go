package tui

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/lipgloss"
)

func (m Model) renderCommitList(width int, maxLines int) string {
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

	end := m.commitOffset + maxLines
	if end > len(m.commits) {
		end = len(m.commits)
	}

	var items []string
	// Available subject width inside panel: totalWidth - padding(2) - hash(7) - gap(2) - prefix(2)
	maxSubject := width - 13
	if maxSubject < 10 {
		maxSubject = 10
	}
	for i := m.commitOffset; i < end; i++ {
		commit := m.commits[i]
		subject := commit.Subject
		if len(subject) > maxSubject {
			subject = subject[:maxSubject-3] + "..."
		}
		line := fmt.Sprintf("%s  %s", StyleMuted.Render(commit.Hash), subject)
		prefix := "  "
		style := StyleNormal
		if i == m.cursor {
			prefix = "▸ "
			style = StyleSelected
		}
		items = append(items, prefix+style.Render(line))
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
