package tui

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/lipgloss"
)

type Branch struct {
	Name    string
	Current bool
	Remote  bool
}

func parseBranches(output string) []Branch {
	lines := strings.Split(strings.TrimSpace(output), "\n")
	branches := make([]Branch, 0, len(lines))
	for _, line := range lines {
		current := strings.HasPrefix(line, "*")
		name := strings.TrimLeft(line, "* ")
		remote := strings.HasPrefix(name, "remotes/")
		if remote {
			name = strings.TrimPrefix(name, "remotes/")
		}
		branches = append(branches, Branch{Name: name, Current: current, Remote: remote})
	}
	return branches
}

func (m Model) renderBranchModal() string {
	modalWidth := 50
	if m.width-4 < modalWidth {
		modalWidth = m.width - 4
	}

	title := StyleAccent.Render("Branches")

	var items []string
	for i, b := range m.branches {
		var line string
		if b.Remote {
			line = "  " + StyleMuted.Render(b.Name)
		} else if i == m.branchCursor {
			prefix := "▸ "
			name := b.Name
			if b.Current {
				name = fmt.Sprintf("%s %s", StyleBold.Render("◉"), name)
			}
			line = StyleSelected.Render(prefix + name)
		} else {
			name := b.Name
			if b.Current {
				name = fmt.Sprintf("  %s %s", StyleBold.Render("◉"), name)
			}
			line = "  " + StyleNormal.Render(name)
		}
		items = append(items, line)
	}
	list := strings.Join(items, "\n")

	footer := StyleMuted.Render("↑↓ nav  enter/e checkout  esc close")

	content := lipgloss.JoinVertical(
		lipgloss.Left,
		title,
		strings.Repeat("─", Max(8, modalWidth-2)),
		list,
		"",
		footer,
	)

	modal := lipgloss.NewStyle().
		Width(modalWidth).
		Border(lipgloss.RoundedBorder()).
		BorderForeground(ColorBorder).
		Padding(0, 1).
		Render(content)

	return lipgloss.Place(
		m.width,
		m.height,
		lipgloss.Center,
		lipgloss.Center,
		modal,
	)
}
