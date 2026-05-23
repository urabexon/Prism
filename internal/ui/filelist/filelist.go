package filelist

import (
	"fmt"

	tea	"github.com/charmbracelet/bubbletea"
	"github.com/urabexon/prism/internal/ghclient"
)

type Model struct {
	prNumber int
	files    []ghclient.File
	cursor   int
}

func New(prNumber int, files []ghclient.File) Model {
	return Model{prNumber: prNumber, files: files}
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
		case "j", "down":
			if m.cursor < len(m.files)-1 {
				m.cursor++
			}
		case "k", "up":
			if m.cursor > 0 {
				m.cursor--
			}
		case "esc":
			// PRリストに戻る
		}
	}
	return m, nil
}

func (m Model) View() string {
	s := fmt.Sprintf("PR #%d - Files\n\n", m.prNumber)
	for i, f := range m.files {
		cursor := " "
		if i == m.cursor {
			cursor = ">"
		}
		s += fmt.Sprintf("%s %s +%d -%d\n", cursor, f.Path, f.Additions, f.Deletions)
	}
	s += "\nj/k: move, esc: back, q: quit"
	return s
}
