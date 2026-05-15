package tui

import (
	"os"
	"os/exec"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

// Model owns UI state for the initial scaffold.
type Model struct {
	width         int
	height        int
	path          string
	repoDetected  bool
	view          string
	status        string
	commits       []string
	detailContent string
	cursor        int
	focus         string
}

// NewModel builds a default model with lightweight runtime context.
func NewModel() Model {
	cwd, err := os.Getwd()
	if err != nil {
		cwd = "unknown"
	}

	_, statErr := os.Stat(".git")
	repoDetected := statErr == nil
	status := "Status: waiting for repository data"
	commits := []string{"- (placeholder)"}

	if !repoDetected {
		status = "Status: no repository detected"
	} else {
		log, err := exec.Command("git", "log", "--oneline", "-20").Output()
		if err == nil {
			out := strings.TrimSpace(string(log))
			if out != "" {
				commits = strings.Split(out, "\n")
				status = "Status: loaded commits"
			}
		}
	}

	return Model{
		path:         cwd,
		repoDetected: repoDetected,
		view:         "commits",
		status:       status,
		commits:      commits,
		focus:        "commits",
	}
}

func (m Model) Init() tea.Cmd {
	return nil
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "q", "ctrl+c":
			return m, tea.Quit
		case "up", "k", "w":
			if m.focus == "commits" && m.cursor > 0 {
				m.cursor--
			}
		case "down", "j", "s":
			if m.focus == "commits" && m.cursor < len(m.commits)-1 {
				m.cursor++
			}
		case "tab":
			if m.focus == "commits" {
				m.focus = "details"
			} else {
				m.focus = "commits"
			}
		case "shift+tab":
			if m.focus == "details" {
				m.focus = "commits"
			} else {
				m.focus = "details"
			}
		case "enter":
			if m.focus == "commits" && m.repoDetected && len(m.commits) > 0 && m.commits[0] != "- (placeholder)" {
				hash := strings.SplitN(m.commits[m.cursor], " ", 2)[0]
				out, err := exec.Command("git", "show", "--stat", hash).Output()
				if err == nil {
					m.detailContent = string(out)
					m.focus = "details"
					m.status = "Status: inspected " + hash
				}
			}
		case "esc":
			if m.focus == "details" {
				m.focus = "commits"
			}
		case "r":
			if m.repoDetected {
				log, err := exec.Command("git", "log", "--oneline", "-20").Output()
				if err == nil {
					out := strings.TrimSpace(string(log))
					if out != "" {
						m.commits = strings.Split(out, "\n")
					}
				}
				if m.cursor >= len(m.commits) {
					m.cursor = len(m.commits) - 1
				}
				m.status = "Status: refreshed"
			}
		case "b":
			m.view = "branches"
			m.status = "Status: branch view (placeholder)"
		}
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
	}

	return m, nil
}

func (m Model) View() string {
	if m.width == 0 {
		m.width = 100
	}
	if m.height == 0 {
		m.height = 30
	}

	frame := FrameStyle.Width(m.width - 2)

	contentWidth := frame.GetWidth() - frame.GetHorizontalPadding()
	if contentWidth < 40 {
		contentWidth = 40
	}

	header := m.renderHeader(contentWidth)
	main := m.renderMain(contentWidth)
	status := m.renderStatus(contentWidth)
	footer := m.renderFooter(contentWidth)

	sections := []string{header, main}
	if m.height >= 8 {
		sections = append(sections, status)
	}
	if m.height >= 6 {
		sections = append(sections, footer)
	}
	body := lipgloss.JoinVertical(lipgloss.Left, sections...)
	return frame.Render(body) + "\n"
}

func Run() error {
	p := tea.NewProgram(NewModel(), tea.WithAltScreen())
	_, err := p.Run()
	return err
}

func (m Model) renderMain(width int) string {
	colGap := 1
	minPaneWidth := 20

	if width >= minPaneWidth*2+colGap {
		leftWidth := (width - colGap) / 2
		rightWidth := width - colGap - leftWidth
		left := m.renderCommitList(leftWidth)
		right := m.renderDetails(rightWidth)
		panes := lipgloss.JoinHorizontal(lipgloss.Top, left, strings.Repeat(" ", colGap), right)
		line := StyleMuted.Render(strings.Repeat("─", width))
		return lipgloss.JoinVertical(lipgloss.Left, panes, line)
	}

	left := m.renderCommitList(width)
	right := m.renderDetails(width)
	sep := StyleMuted.Render(strings.Repeat("─", width))
	return lipgloss.JoinVertical(lipgloss.Left, left, sep, right)
}
