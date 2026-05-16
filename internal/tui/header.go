package tui

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/lipgloss"
)

func (m Model) renderHeader(width int) string {
	title := StyleAccent.Render("gitviz")

	branchInfo := "no repo"
	if m.repoDetected && m.currentBranch != "" {
		branchInfo = m.currentBranch
	} else if m.repoDetected {
		branchInfo = "detached HEAD"
	}

	meta := fmt.Sprintf("%s  |  View: %s", StyleAccent.Render(branchInfo), m.view)
	path := fmt.Sprintf("Path: %s", m.path)
	line := strings.Repeat("─", width)

	return lipgloss.JoinVertical(
		lipgloss.Left,
		fmt.Sprintf("%s  %s", title, StyleMuted.Render(meta)),
		StyleMuted.Render(path),
		StyleMuted.Render(line),
	)
}
