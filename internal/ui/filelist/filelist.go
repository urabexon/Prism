package filelist

import (
	"fmt"
	"strings"

	tea "github.com/charmbracelet/bubbletea"

	"github.com/urabexon/prism/internal/ghclient"
	"github.com/urabexon/prism/internal/state"
	"github.com/urabexon/prism/internal/ui/checks"
	"github.com/urabexon/prism/internal/ui/styles"
)

type SelectFileMsg struct {
	Index int
}

type BackMsg struct{}

type OpenChecksMsg struct {
	PR ghclient.PR
}

type MergeMsg struct {
	Number  int
	Method  string
	Undraft bool
}

type Model struct {
	pr             ghclient.PR
	diff           ghclient.ParsedDiff
	cursor         int
	width          int
	height         int
	repo           string
	store          *state.Store
	loading        bool
	err            error
	confirmMerge   bool
	mergeMethod    int
	merging        bool
	mergeResult    string
	allowedMethods []string
}

func New(repo string, store *state.Store) Model {
	return Model{
		repo:  repo,
		store: store,
	}
}

func (m Model) SetPR(pr ghclient.PR) Model {
	m.pr = pr
	m.loading = true
	m.cursor = 0
	return m
}

func (m Model) SetDiff(diff ghclient.ParsedDiff) Model {
	m.diff = diff
	m.loading = false
	return m
}

func (m Model) SetAllowedMergeMethods(methods []string) Model {
	m.allowedMethods = methods
	return m
}

func (m Model) mergeMethods() []string {
	if len(m.allowedMethods) > 0 {
		return m.allowedMethods
	}
	return []string{"squash", "merge", "rebase"}
}

func (m Model) SetMergeResult(msg string) Model {
	m.mergeResult = msg
	m.merging = false
	return m
}

func (m Model) SetError(err error) Model {
	m.err = err
	m.loading = false
	return m
}

func (m Model) SetSize(w, h int) Model {
	m.width = w
	m.height = h
	return m
}

func (m Model) FileCount() int {
	return len(m.diff.Files)
}

func (m Model) Diff() ghclient.ParsedDiff {
	return m.diff
}

func (m Model) PR() ghclient.PR {
	return m.pr
}

func (m Model) visibleHeight() int {
	h := m.height - 7
	if h < 1 {
		return 10
	}
	return h
}

func (m Model) Update(msg tea.Msg) (Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		if m.mergeResult != "" {
			m.mergeResult = ""
		}

		if m.confirmMerge {
			switch msg.String() {
			case "h", "left":
				m.mergeMethod = (m.mergeMethod + len(m.mergeMethods()) - 1) % len(m.mergeMethods())
			case "l", "right":
				m.mergeMethod = (m.mergeMethod + 1) % len(m.mergeMethods())
			case "enter", "y":
				number := m.pr.Number
				method := m.mergeMethods()[m.mergeMethod]
				undraft := m.pr.IsDraft
				m.confirmMerge = false
				m.merging = true
				return m, func() tea.Msg {
					return MergeMsg{Number: number, Method: method, Undraft: undraft}
				}
			case "esc", "n", "q":
				m.confirmMerge = false
			}
			return m, nil
		}

		switch msg.String() {
		case "j", "down", "ctrl+n":
			if m.cursor < len(m.diff.Files)-1 {
				m.cursor++
			}
		case "k", "up", "ctrl+p":
			if m.cursor > 0 {
				m.cursor--
			}
		case "ctrl+d":
			m.cursor = min(m.cursor+m.visibleHeight()/2, len(m.diff.Files)-1)
		case "ctrl+u":
			m.cursor = max(0, m.cursor-m.visibleHeight()/2)
		case "ctrl+f", "pgdown":
			m.cursor = min(m.cursor+m.visibleHeight(), len(m.diff.Files)-1)
		case "ctrl+b", "pgup":
			m.cursor = max(0, m.cursor-m.visibleHeight())
		case "g", "home":
			m.cursor = 0
		case "G", "end":
			if len(m.diff.Files) > 0 {
				m.cursor = len(m.diff.Files) - 1
			}
		case "enter":
			if len(m.diff.Files) > 0 {
				idx := m.cursor
				return m, func() tea.Msg { return SelectFileMsg{Index: idx} }
			}
		case "space", " ":
			if len(m.diff.Files) > 0 {
				file := m.diff.Files[m.cursor]
				m.store.ToggleFileReviewed(m.repo, m.pr.Number, file.FilePath())
				_ = m.store.Save()
			}
		case "a":
			paths := make([]string, len(m.diff.Files))
			for i, f := range m.diff.Files {
				paths[i] = f.FilePath()
			}
			m.store.MarkAllReviewed(m.repo, m.pr.Number, paths)
			_ = m.store.Save()
		case "c":
			pr := m.pr
			return m, func() tea.Msg { return OpenChecksMsg{PR: pr} }
		case "m":
			if !m.merging {
				m.confirmMerge = true
				m.mergeMethod = 0
			}
		case "esc", "backspace":
			return m, func() tea.Msg { return BackMsg{} }
		}
	}
	return m, nil
}

func (m Model) View() string {
	var b strings.Builder

	prTitle := fmt.Sprintf(" #%d %s ", m.pr.Number, m.pr.Title)
	b.WriteString(styles.Title.Render(prTitle))
	b.WriteString("\n")
	authorLine := fmt.Sprintf("  %s → %s  by %s",
		m.pr.HeadRef, m.pr.BaseRef, m.pr.Author)
	b.WriteString(styles.Subtitle.Render(authorLine))
	b.WriteString("\n")
	checkLine := checks.CheckSummaryLine(m.pr.CheckSummary)
	if checkLine != "" {
		b.WriteString(checkLine)
	}
	b.WriteString("\n")

	if m.loading {
		b.WriteString("\n  Loading diff...")
		return b.String()
	}

	if m.err != nil {
		b.WriteString(fmt.Sprintf("\n  Error: %v", m.err))
		return b.String()
	}

	fileCount := len(m.diff.Files)
	reviewed := m.store.ReviewedFileCount(m.repo, m.pr.Number)
	progress := fmt.Sprintf("  Files: %d/%d reviewed", reviewed, fileCount)
	b.WriteString(styles.Subtitle.Render(progress))
	b.WriteString("\n\n")
	visibleHeight := m.height - 7
	if visibleHeight < 1 {
		visibleHeight = 10
	}

	startIdx := 0
	if m.cursor >= visibleHeight {
		startIdx = m.cursor - visibleHeight + 1
	}
	endIdx := startIdx + visibleHeight
	if endIdx > fileCount {
		endIdx = fileCount
	}

	for i := startIdx; i < endIdx; i++ {
		file := m.diff.Files[i]
		selected := i == m.cursor
		line := m.renderFile(file, selected)
		b.WriteString(line)
		b.WriteString("\n")
	}

	if m.merging {
		b.WriteString("\n")
		b.WriteString(styles.Subtitle.Render("  Merging..."))
		b.WriteString("\n")
	} else if m.mergeResult != "" {
		b.WriteString("\n")
		b.WriteString(fmt.Sprintf("  %s", m.mergeResult))
		b.WriteString("\n")
	} else if m.confirmMerge {
		b.WriteString("\n")
		action := "Merge"
		if m.pr.IsDraft {
			action = "Undraft & Merge"
		}
		b.WriteString(styles.Unread.Render(fmt.Sprintf("  %s #%d? ", action, m.pr.Number)))
		for i, method := range m.mergeMethods() {
			if i == m.mergeMethod {
				b.WriteString(styles.Selected.Render(fmt.Sprintf(" [%s] ", method)))
			} else {
				b.WriteString(styles.Help.Render(fmt.Sprintf("  %s  ", method)))
			}
		}
		b.WriteString(styles.Help.Render("  h/l:method  enter:confirm  esc:cancel"))
		b.WriteString("\n")
	}

	b.WriteString("\n")
	help := styles.Help.Render("  j/k:navigate  enter:diff  space:reviewed  c:checks  m:merge  a:all  esc:back")
	b.WriteString(help)

	return b.String()
}

func (m Model) renderFile(file ghclient.FileDiff, selected bool) string {
	path := file.FilePath()
	isReviewed := m.store.IsFileReviewed(m.repo, m.pr.Number, path)
	indicator := "  "
	if isReviewed {
		indicator = styles.Reviewed.Render("✓ ")
	}

	status := ""
	switch {
	case file.IsNew:
		status = styles.Added.Render("[new] ")
	case file.IsDelete:
		status = styles.Removed.Render("[del] ")
	case file.IsRename:
		status = styles.Subtitle.Render("[ren] ")
	case file.IsBinary:
		status = styles.Subtitle.Render("[bin] ")
	}

	adds := file.Additions()
	dels := file.Deletions()
	stats := fmt.Sprintf("%s%s",
		styles.Added.Render(fmt.Sprintf("+%-4d", adds)),
		styles.Removed.Render(fmt.Sprintf("-%-4d", dels)),
	)
	pathWidth := m.width - 30
	if pathWidth < 20 {
		pathWidth = 20
	}
	if len(path) > pathWidth {
		path = "…" + path[len(path)-pathWidth+1:]
	}

	line := fmt.Sprintf("%s%s%s  %s", indicator, status, path, stats)

	if selected {
		line = styles.Selected.Width(m.width).Render(line)
	}

	return line
}
