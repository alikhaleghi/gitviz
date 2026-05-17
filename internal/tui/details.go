package tui

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/lipgloss"
)

func (m Model) renderDetails(width int, detailHeight int) string {
	titleStyle := StyleBold
	if m.focus == "details" {
		titleStyle = StyleAccent
	}

	var body string
	switch {
	case m.detail != nil:
		body = m.formatDetail(width, detailHeight)
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

func (m Model) formatDetail(width int, detailHeight int) string {
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

	nonDiffLines := 4

	if d.Body != "" {
		body := lipgloss.NewStyle().Width(inner).Render(d.Body)
		sections = append(sections, "", body, sep)
		nonDiffLines += 3
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
		nonDiffLines += 2 + len(d.Files)
	}

	if d.Diff != "" {
		diffHeader := StyleBold.Render("── Diff ──")
		allLines := strings.Split(d.Diff, "\n")
		var colored []string
		for _, line := range allLines {
			if len(line) > 0 {
				switch line[0] {
				case '+':
					colored = append(colored, StyleAdd.Render(line))
				case '-':
					colored = append(colored, StyleDel.Render(line))
				case '@':
					colored = append(colored, StyleHunk.Render(line))
				default:
					colored = append(colored, StyleMuted.Render(line))
				}
			} else {
				colored = append(colored, "")
			}
		}

		overhead := 4
		if d.Body != "" {
			overhead++
		}
		diffBudget := Max(0, detailHeight-nonDiffLines-overhead)
		if m.detailFocus == "diff" {
			diffBudget = Max(0, detailHeight-nonDiffLines-1)
		}

		start := m.diffScroll
		if start >= len(colored) {
			start = Max(0, len(colored)-1)
		}
		end := Min(start+diffBudget, len(colored))
		if end <= start {
			end = Min(start+1, len(colored))
		}
		visible := colored[start:end]

		prefix := ""
		if start > 0 {
			prefix = StyleMuted.Render("  ... (scroll with j/k)")
		}

		diffLines := append([]string{prefix}, visible...)
		sections = append(sections, "", diffHeader, strings.Join(diffLines, "\n"))
	}

	return lipgloss.JoinVertical(lipgloss.Left, sections...)
}
