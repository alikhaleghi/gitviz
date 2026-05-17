package tui

import (
	"fmt"
	"os"
	"os/exec"
	"runtime"
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
	Diff   string
}

// Model owns UI state for the initial scaffold.
type Model struct {
	width          int
	height         int
	path           string
	repoDetected   bool
	view           string
	status         string
	commits        []Commit
	detail         *CommitDetail
	cursor         int
	focus          string
	showModal      bool
	branches       []Branch
	branchCursor   int
	currentBranch  string
	commitOffset   int
	commitMaxLines int
	detailFocus    string
	diffScroll     int
	blameView      bool
	blameFile      string
	blameLines     []string
	blameScroll    int
	branchMsg      string
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
			if m.showModal {
				m.branchMsg = ""
				for i := m.branchCursor - 1; i >= 0; i-- {
					if !m.branches[i].Remote {
						m.branchCursor = i
						break
					}
				}
			} else if m.blameView {
				if m.blameScroll > 0 {
					m.blameScroll--
				}
			} else if m.detailFocus == "diff" {
				if m.diffScroll > 0 {
					m.diffScroll--
				}
			} else if m.focus == "commits" && m.cursor > 0 {
				m.cursor--
				if m.cursor < m.commitOffset {
					m.commitOffset = m.cursor
				}
				m = m.loadDetail(m.cursor)
			}
		case "down", "j", "s":
			if m.showModal {
				m.branchMsg = ""
				for i := m.branchCursor + 1; i < len(m.branches); i++ {
					if !m.branches[i].Remote {
						m.branchCursor = i
						break
					}
				}
			} else if m.blameView {
				if m.blameScroll < len(m.blameLines)-1 {
					m.blameScroll++
				}
			} else if m.detailFocus == "diff" {
				diffLines := strings.Count(m.detail.Diff, "\n") + 1
				if m.diffScroll < diffLines-1 {
					m.diffScroll++
				}
			} else if m.focus == "commits" && m.cursor < len(m.commits)-1 {
				m.cursor++
				if m.commitMaxLines > 0 && m.cursor >= m.commitOffset+m.commitMaxLines {
					m.commitOffset = m.cursor - m.commitMaxLines + 1
				}
				m = m.loadDetail(m.cursor)
			}
		case "tab":
			if m.showModal {
				m.showModal = false
				m.branchMsg = ""
			} else if m.focus == "commits" {
				m.focus = "details"
			} else if m.detailFocus == "diff" {
				m.detailFocus = ""
			} else if m.detail != nil && m.detail.Diff != "" {
				m.detailFocus = "diff"
			}
		case "shift+tab":
			if m.showModal {
				m.showModal = false
				m.branchMsg = ""
			} else if m.focus == "details" && m.detailFocus == "diff" {
				m.detailFocus = ""
			} else if m.focus == "details" {
				m.focus = "commits"
			} else {
				m.focus = "details"
			}
		case "enter", "e":
			if m.showModal && m.branchCursor < len(m.branches) {
				b := m.branches[m.branchCursor]
				if b.Remote {
					m.branchMsg = "Cannot checkout a remote branch"
				} else if b.Current {
					m.branchMsg = fmt.Sprintf("Already on %s", b.Name)
				} else {
					out, err := exec.Command("git", "checkout", b.Name).CombinedOutput()
					if err != nil {
						msg := strings.TrimSpace(string(out))
						if msg == "" {
							msg = err.Error()
						}
						firstLine, _, _ := strings.Cut(msg, "\n")
						if len(firstLine) > 60 {
							firstLine = firstLine[:60] + "..."
						}
						m.branchMsg = firstLine
					} else {
						m.currentBranch = b.Name
						m.showModal = false
						m.branchMsg = ""
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
		case "B":
			if m.detail != nil && !m.blameView {
				if len(m.detail.Files) == 0 {
					m.status = "Status: no files to blame"
				} else {
					path := m.detail.Files[0]
					if parts := strings.SplitN(path, "\t", 2); len(parts) == 2 {
						path = parts[1]
					}
					m = m.loadBlame(path)
				}
			}
		case "esc":
			if m.showModal {
				m.showModal = false
				m.branchMsg = ""
			} else if m.blameView {
				m.blameView = false
				m.blameFile = ""
				m.blameLines = nil
				m.blameScroll = 0
				m.focus = "details"
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
				m.branchMsg = ""
				out, err := exec.Command("git", "branch", "--all", "--no-color").Output()
				if err == nil {
					m.branches = parseBranches(string(out))
				}
				m.status = fmt.Sprintf("Status: %d branches", len(m.branches))
			} else if m.showModal {
				m.showModal = false
				m.branchMsg = ""
			}
		case "y":
			if m.detail != nil {
				hash := m.commits[m.cursor].Hash
				if err := copyToClipboard(hash); err != nil {
					m.status = fmt.Sprintf("Failed to copy: %s", err)
				} else {
					short := hash
					if len(short) > 8 {
						short = hash[:8] + "..."
					}
					m.status = fmt.Sprintf("Copied: %s", short)
				}
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

	// The frame adds rounded border (2 wide, top + bottom lines)
	frame := FrameStyle.Width(m.width - 2)
	contentWidth := frame.GetWidth() - frame.GetHorizontalPadding()
	if contentWidth < 40 {
		contentWidth = 40
	}

	// Compute vertical space for the main pane row
	// Outside renderMain: frame border (2) + header (3)
	mainHeight := m.height - 2 - 3
	if m.height >= 8 {
		mainHeight -= 2 // status: text + separator
	}
	if m.height >= 6 {
		mainHeight -= 1 // footer: help text
	}
	if mainHeight < 3 {
		mainHeight = 3
	}

	header := m.renderHeader(contentWidth)
	main := m.renderMain(contentWidth, mainHeight)
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
	rendered := frame.Render(body)

	if m.showModal {
		rendered = m.renderBranchModal()
	}

	lines := strings.Split(rendered, "\n")
	if len(lines) > m.height {
		lines = lines[:m.height]
	}
	return strings.Join(lines, "\n")
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

	diff, err := exec.Command("git", "show", "--format=", "--unified=3", hash).Output()
	if err == nil {
		detail.Diff = strings.TrimSuffix(string(diff), "\n")
	}

	return detail, nil
}

func (m Model) clampScroll(maxLines int) Model {
	if m.cursor < m.commitOffset {
		m.commitOffset = m.cursor
	}
	if m.commitOffset > m.cursor {
		m.commitOffset = m.cursor
	}
	if maxLines > 0 && m.commitOffset > len(m.commits)-maxLines {
		m.commitOffset = len(m.commits) - maxLines
	}
	if m.commitOffset < 0 {
		m.commitOffset = 0
	}
	if len(m.commits) > 0 && m.cursor >= m.commitOffset+maxLines {
		m.commitOffset = m.cursor - maxLines + 1
	}
	return m
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

func (m Model) loadBlame(filePath string) Model {
	hash := m.commits[m.cursor].Hash
	out, err := exec.Command("git", "blame", "--date=short", hash, "--", filePath).Output()
	if err != nil {
		m.status = fmt.Sprintf("Failed to blame: %s", err)
		return m
	}
	m.blameView = true
	m.blameFile = filePath
	m.blameLines = strings.Split(strings.TrimSuffix(string(out), "\n"), "\n")
	m.blameScroll = 0
	m.focus = "details"
	m.status = fmt.Sprintf("Blame: %s (%d lines)", filePath, len(m.blameLines))
	return m
}

func copyToClipboard(text string) error {
	var cmd *exec.Cmd
	switch runtime.GOOS {
	case "darwin":
		cmd = exec.Command("pbcopy")
	case "windows":
		cmd = exec.Command("clip")
	default:
		if _, err := exec.LookPath("xclip"); err == nil {
			cmd = exec.Command("xclip", "-selection", "clipboard")
		} else if _, err := exec.LookPath("xsel"); err == nil {
			cmd = exec.Command("xsel", "--clipboard", "--input")
		} else {
			return fmt.Errorf("no clipboard command found (install xclip or xsel)")
		}
	}
	cmd.Stdin = strings.NewReader(text)
	return cmd.Run()
}

func Run() error {
	p := tea.NewProgram(NewModel(), tea.WithAltScreen())
	_, err := p.Run()
	return err
}

func (m Model) renderMain(width int, mainHeight int) string {
	colGap := 1
	minPaneWidth := 20

	if width >= minPaneWidth*2+colGap {
		leftWidth := (width - colGap) / 2
		rightWidth := width - colGap - leftWidth

		m.commitMaxLines = Max(1, mainHeight-3)
		m = m.clampScroll(m.commitMaxLines)
		left := m.renderCommitList(leftWidth, m.commitMaxLines)
		right := m.renderDetails(rightWidth, mainHeight)

		rightLines := strings.Split(right, "\n")
		if len(rightLines) > mainHeight-1 {
			rightLines = rightLines[:mainHeight-1]
		}
		right = strings.Join(rightLines, "\n")

		panes := lipgloss.JoinHorizontal(lipgloss.Top, left, strings.Repeat(" ", colGap), right)
		sep := StyleMuted.Render(strings.Repeat("─", width))
		combined := lipgloss.JoinVertical(lipgloss.Left, panes, sep)
		combinedLines := strings.Split(combined, "\n")
		if len(combinedLines) > mainHeight {
			combinedLines = combinedLines[:mainHeight]
		}
		return strings.Join(combinedLines, "\n")
	}

	m.commitMaxLines = Max(1, (mainHeight-2)/2)
	m = m.clampScroll(m.commitMaxLines)
	left := m.renderCommitList(width, m.commitMaxLines)
	right := m.renderDetails(width, mainHeight)
	sep := StyleMuted.Render(strings.Repeat("─", width))

	combined := lipgloss.JoinVertical(lipgloss.Left, left, sep, right)
	combinedLines := strings.Split(combined, "\n")
	if len(combinedLines) > mainHeight {
		combinedLines = combinedLines[:mainHeight]
	}
	return strings.Join(combinedLines, "\n")
}
