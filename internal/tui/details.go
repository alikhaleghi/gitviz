package tui

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/lipgloss"
)

func (m Model) renderDetails(width int) string {
	titleStyle := StyleBold
	if m.focus == "details" {
		titleStyle = StyleAccent
	}

	var body string
	switch {
	case m.detail != nil:
		body = m.formatDetail(width)
	case !m.repoDetected:
		body = "Welcome to gitviz\n\nTo get started here:\n1. git init\n2. git add .\n3. git commit -m \"chore: initialize project scaffold\""
	default:
		body = "Select a commit to inspect."
	}

	return StylePanel.Width(width).Render(
		lipgloss.JoinVertical(
			lipgloss.Left,
			titleStyle.Render("Details"),
			strings.Repeat("─", Max(8, width-2)),
			body,
		),
	)
}

func (m Model) formatDetail(width int) string {
	d := m.detail
	inner := width - 2
	if inner < 20 {
		inner = 20
	}

	hash := StyleMuted.Render(d.Hash)
	author := fmt.Sprintf("Author:  %s  <%s>", d.Author, d.Email)
	date := fmt.Sprintf("Date:    %s", d.Date)

	sep := strings.Repeat("─", inner)

	var sections []string
	sections = append(sections, hash, author, date)

	if d.Body != "" {
		body := lipgloss.NewStyle().Width(inner).Render(d.Body)
		sections = append(sections, "", body, sep)
	}

	if len(d.Files) > 0 {
		header := StyleBold.Render("── Files ──")
		var fileLines []string
		for _, f := range d.Files {
			parts := strings.SplitN(f, "\t", 2)
			if len(parts) == 2 {
				fileLines = append(fileLines, fmt.Sprintf("%s  %s", parts[0], parts[1]))
			} else {
				fileLines = append(fileLines, f)
			}
		}
		sections = append(sections, header, strings.Join(fileLines, "\n"))
	}

	return lipgloss.JoinVertical(lipgloss.Left, sections...)
}
