package tui

import (
	"fmt"
	"os"
	"os/exec"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

// Commit holds structured data from git log.
type Commit struct {
	Hash    string
	Subject string
}

// CommitDetail holds parsed output for a single inspected commit.
type CommitDetail struct {
	Hash   string
	Author string
	Email  string
	Date   string
	Body   string
	Files  []string
}

// Model owns UI state for the initial scaffold.
type Model struct {
	width         int
	height        int
	path          string
	repoDetected  bool
	view          string
	status        string
	commits       []Commit
	detail        *CommitDetail
	cursor        int
	focus         string
	showModal     bool
	branches      []Branch
	branchCursor  int
	currentBranch string
}

// NewModel builds a default model with lightweight runtime context.
func NewModel() Model {
	cwd, err := os.Getwd()
	if err != nil {
		cwd = "unknown"
	}

	_, statErr := os.Stat(".git")
	repoDetected := statErr == nil
	status := "Status: no repository detected"
	commits := []Commit{}
	currentBranch := ""

	if repoDetected {
		log, err := exec.Command("git", "log", "--oneline", "-n", "50").Output()
		if err == nil {
			commits = parseGitLog(string(log))
			status = fmt.Sprintf("Status: loaded %d commits", len(commits))
		}
		branch, _ := exec.Command("git", "branch", "--show-current").Output()
		currentBranch = strings.TrimSpace(string(branch))
	}

	m := Model{
		path:          cwd,
		repoDetected:  repoDetected,
		view:          "commits",
		status:        status,
		commits:       commits,
		focus:         "commits",
		currentBranch: currentBranch,
	}
	if len(commits) > 0 {
		m = m.loadDetail(0)
	}
	return m
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
			if m.showModal && m.branchCursor > 0 {
				m.branchCursor--
			} else if m.focus == "commits" && m.cursor > 0 {
				m.cursor--
				m = m.loadDetail(m.cursor)
			}
		case "down", "j", "s":
			if m.showModal && m.branchCursor < len(m.branches)-1 {
				m.branchCursor++
			} else if m.focus == "commits" && m.cursor < len(m.commits)-1 {
				m.cursor++
				m = m.loadDetail(m.cursor)
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
		case "enter", "e":
			if m.showModal && m.branchCursor < len(m.branches) {
				b := m.branches[m.branchCursor]
				if b.Remote {
					m.status = "Cannot checkout a remote branch"
				} else if b.Current {
					m.showModal = false
					m.status = fmt.Sprintf("Already on %s", b.Name)
				} else {
					err := exec.Command("git", "checkout", b.Name).Run()
					if err != nil {
						m.status = fmt.Sprintf("Failed to checkout %s", b.Name)
					} else {
						m.currentBranch = b.Name
						m.showModal = false
						m.status = fmt.Sprintf("Switched to %s", b.Name)
						log, logErr := exec.Command("git", "log", "--oneline", "-n", "50").Output()
						if logErr == nil {
							m.commits = parseGitLog(string(log))
						}
						if m.cursor >= len(m.commits) {
							m.cursor = len(m.commits) - 1
						}
						if len(m.commits) > 0 {
							m = m.loadDetail(m.cursor)
						}
					}
				}
			} else if m.focus == "commits" && m.cursor < len(m.commits) {
				m = m.loadDetail(m.cursor)
				if m.detail != nil {
					m.focus = "details"
				}
			}
		case "esc":
			if m.showModal {
				m.showModal = false
			} else if m.focus == "details" {
				m.focus = "commits"
			}
		case "r":
			if m.repoDetected {
				log, err := exec.Command("git", "log", "--oneline", "-n", "50").Output()
				if err == nil {
					m.commits = parseGitLog(string(log))
				}
				if m.cursor >= len(m.commits) {
					m.cursor = len(m.commits) - 1
				}
				m.status = "Status: refreshed"
			}
		case "b":
			if m.repoDetected && !m.showModal {
				m.showModal = true
				m.branchCursor = 0
				out, err := exec.Command("git", "branch", "--all", "--no-color").Output()
				if err == nil {
					m.branches = parseBranches(string(out))
				}
				m.status = fmt.Sprintf("Status: %d branches", len(m.branches))
			} else if m.showModal {
				m.showModal = false
			}
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
	rendered := frame.Render(body) + "\n"

	if m.showModal {
		rendered = m.renderBranchModal()
	}

	return rendered
}

func parseGitLog(output string) []Commit {
	lines := strings.Split(strings.TrimSpace(output), "\n")
	commits := make([]Commit, 0, len(lines))
	for _, line := range lines {
		parts := strings.SplitN(line, " ", 2)
		if len(parts) == 2 {
			commits = append(commits, Commit{Hash: parts[0], Subject: parts[1]})
		}
	}
	return commits
}

func fetchCommitDetail(hash string) (CommitDetail, error) {
	out, err := exec.Command("git", "log", "-1", "--format=%H%n%an%n%ae%n%aI%n%b", hash).Output()
	if err != nil {
		return CommitDetail{}, err
	}
	lines := strings.SplitN(string(out), "\n", 5)
	if len(lines) < 4 {
		return CommitDetail{}, fmt.Errorf("unexpected git log output")
	}

	detail := CommitDetail{
		Hash:   lines[0],
		Author: lines[1],
		Email:  lines[2],
		Date:   lines[3],
	}
	if len(lines) > 4 {
		detail.Body = strings.TrimSpace(lines[4])
	}

	files, err := exec.Command("git", "diff-tree", "--no-commit-id", "-r", "--name-status", hash).Output()
	if err == nil {
		raw := strings.TrimSpace(string(files))
		if raw != "" {
			detail.Files = strings.Split(raw, "\n")
		}
	}

	return detail, nil
}

func (m Model) loadDetail(idx int) Model {
	if !m.repoDetected || idx >= len(m.commits) {
		m.detail = nil
		return m
	}
	hash := m.commits[idx].Hash
	detail, err := fetchCommitDetail(hash)
	if err != nil {
		m.detail = nil
		m.status = "Status: failed to load commit details"
		return m
	}
	m.detail = &detail
	m.status = fmt.Sprintf("Status: inspected %s", hash)
	return m
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
