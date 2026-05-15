package tui

import "github.com/charmbracelet/lipgloss"

var (
	ColorMuted  = lipgloss.Color("8")
	ColorAccent = lipgloss.Color("12")
	ColorNormal = lipgloss.Color("7")
	ColorBorder = lipgloss.Color("8")

	BorderRounded = lipgloss.RoundedBorder()

	StyleMuted  = lipgloss.NewStyle().Foreground(ColorMuted)
	StyleAccent = lipgloss.NewStyle().Bold(true).Foreground(ColorAccent)
	StyleBold   = lipgloss.NewStyle().Bold(true)
	StyleNormal = lipgloss.NewStyle().Foreground(ColorNormal)

	StyleSep   = lipgloss.NewStyle().Foreground(ColorMuted)
	StylePanel = lipgloss.NewStyle().Padding(0, 1)

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
