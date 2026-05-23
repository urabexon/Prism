package main

import (
	"fmt"
	"log"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/urabexon/prism/internal/ghclient"
)

type model struct {
	prs      []ghclient.PR
	cursor   int
	selected map[int]struct{}
}

func initialModel(prs []ghclient.PR) model {
	return model{
		prs:      prs,
		selected: make(map[int]struct{}),
	}
}

func (m model) Init() tea.Cmd {
	return nil
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "q":
			return m, tea.Quit
		case "up", "k":
			if m.cursor > 0 {
				m.cursor--
			}
		case "down", "j":
			if m.cursor < len(m.prs)-1 {
				m.cursor++
			}
		case "enter", " ":
			if _, ok := m.selected[m.cursor]; ok {
				delete(m.selected, m.cursor)
			} else {
				m.selected[m.cursor] = struct{}{}
			}
		}
	}
	return m, nil
}

func (m model) View() string {
	s := "Pull Requests\n\n"

	for i, pr := range m.prs {
		cursor := " "
		if m.cursor == i {
			cursor = ">"
		}

		checked := " "
		if _, ok := m.selected[i]; ok {
			checked = "x"
		}

		s += fmt.Sprintf("%s [%s] #%d %s\n", cursor, checked, pr.Number, pr.Title)
	}

	s += "\nPress q to quit.\n"
	return s
}

func main() {
	prs, err := ghclient.ListPRs()
	if err != nil {
		log.Fatal(err)
	}

	p := tea.NewProgram(initialModel(prs))
	if _, err := p.Run(); err != nil {
		log.Fatal(err)
	}
}
