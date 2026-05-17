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
	case m.compareMode:
		body = m.formatCompare(width, detailHeight)
	case m.blameView:
		body = m.formatBlame(width, detailHeight)
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

func (m Model) formatBlame(width int, detailHeight int) string {
	inner := width - 2
	if inner < 20 {
		inner = 20
	}

	title := StyleBold.Render("Blame: " + m.blameFile)
	fileLines := strings.Split(m.blameFile, "/")
	shortFile := fileLines[len(fileLines)-1]
	budget := Max(0, detailHeight-3)

	start := m.blameScroll
	if start >= len(m.blameLines) {
		start = Max(0, len(m.blameLines)-1)
	}
	end := Min(start+budget, len(m.blameLines))
	if end <= start {
		end = Min(start+1, len(m.blameLines))
	}
	visible := m.blameLines[start:end]

	prefix := ""
	if start > 0 {
		prefix = StyleMuted.Render("  ... (scroll with j/k)")
	}

	lines := append([]string{prefix}, visible...)
	colored := make([]string, 0, len(lines))
	for _, line := range lines {
		if idx := strings.Index(line, ")"); idx != -1 && idx < 30 {
			hash := StyleAccent.Render(line[:7])
			meta := StyleMuted.Render(line[7 : idx+1])
			content := line[idx+1:]
			if len(content) > 0 && content[0] == ' ' {
				content = content[1:]
			}
			colored = append(colored, fmt.Sprintf("%s %s %s", hash, meta, content))
		} else {
			colored = append(colored, StyleMuted.Render(line))
		}
	}

	return lipgloss.JoinVertical(lipgloss.Left,
		title,
		StyleMuted.Render("Esc to return  │  j/k scroll  │  "+shortFile),
		strings.Join(colored, "\n"),
	)
}

func (m Model) formatCompare(width int, detailHeight int) string {
	inner := width - 2
	if inner < 20 {
		inner = 20
	}

	first := m.commits[m.compareSelected[0]].Hash[:8]
	second := m.commits[m.compareSelected[1]].Hash[:8]
	title := StyleBold.Render(fmt.Sprintf("Compare: %s..%s", first, second))
	budget := Max(0, detailHeight-3)

	allLines := strings.Split(m.compareDiff, "\n")
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

	start := m.compareScroll
	if start >= len(colored) {
		start = Max(0, len(colored)-1)
	}
	end := Min(start+budget, len(colored))
	if end <= start {
		end = Min(start+1, len(colored))
	}
	visible := colored[start:end]

	prefix := ""
	if start > 0 {
		prefix = StyleMuted.Render("  ... (scroll with j/k)")
	}

	lines := append([]string{prefix}, visible...)
	return lipgloss.JoinVertical(lipgloss.Left,
		title,
		StyleMuted.Render("Esc to clear compare"),
		strings.Join(lines, "\n"),
	)
}
