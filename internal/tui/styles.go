package tui

import "github.com/charmbracelet/lipgloss"

var (
	ColorMuted  = lipgloss.Color("8")
	ColorAccent = lipgloss.Color("12")
	ColorNormal = lipgloss.Color("7")
	ColorBorder = lipgloss.Color("8")
	ColorAdd    = lipgloss.Color("2")
	ColorDel    = lipgloss.Color("9")
	ColorHunk   = lipgloss.Color("3")

	BorderRounded = lipgloss.RoundedBorder()

	StyleMuted  = lipgloss.NewStyle().Foreground(ColorMuted)
	StyleAccent = lipgloss.NewStyle().Bold(true).Foreground(ColorAccent)
	StyleBold   = lipgloss.NewStyle().Bold(true)
	StyleNormal = lipgloss.NewStyle().Foreground(ColorNormal)
	StyleAdd    = lipgloss.NewStyle().Foreground(ColorAdd)
	StyleDel    = lipgloss.NewStyle().Foreground(ColorDel)
	StyleHunk   = lipgloss.NewStyle().Foreground(ColorHunk)

	StyleSep      = lipgloss.NewStyle().Foreground(ColorMuted)
	StylePanel    = lipgloss.NewStyle().Padding(0, 1)
	StyleSelected = lipgloss.NewStyle().Bold(true).Foreground(ColorAccent)

	FrameStyle = lipgloss.NewStyle().
			Border(BorderRounded).
			BorderForeground(ColorBorder).
			Padding(0, 1)
)

func Max(a, b int) int {
	if a > b {
		return a
	}
	return b
}

func Min(a, b int) int {
	if a < b {
		return a
	}
	return b
}
