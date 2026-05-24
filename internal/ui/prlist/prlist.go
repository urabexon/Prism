package prlist

import (
	"fmt"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/urabexon/prism/internal/ghclient"
)

type Model struct {
	prs    []ghclient.PR
	cursor int
}

func New(prs []ghclient.PR) Model {
	return Model{prs: prs}
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
			if m.cursor < len(m.prs)-1 {
				m.cursor++
			}
		case "k", "up":
			if m.cursor > 0 {
				m.cursor--
			}
		case "enter":
		}
	}
	return m, nil
}

func (m Model) View() string {
	s := "PR List\n\n"
	for i, pr := range m.prs {
		cursor := " "
		if i == m.cursor {
			cursor = ">"
		}
		s += fmt.Sprintf("%s #%d %s\n", cursor, pr.Number, pr.Title)
	}
	s += "\nj/k: move, enter: select, q: quit"
	return s
}

func (m Model) SelectedPR() ghclient.PR {
	return m.prs[m.cursor]
}
