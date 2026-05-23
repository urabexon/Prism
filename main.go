package main

import (
	"fmt"
	"log"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/urabexon/prism/internal/ghclient"
)

type model struct {
	prs        []ghclient.PR
	cursor     int
	showDetail bool
	detail     *ghclient.PRDetail
	showFiles  bool
	files      []ghclient.File
}

func initialModel(prs []ghclient.PR) model {
	return model{
		prs: prs,
	}
}

func (m model) Init() tea.Cmd {
	return nil
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		// ファイル一覧表示中の場合
		if m.showFiles {
			switch msg.String() {
			case "ctrl+c", "q":
				return m, tea.Quit
			case "esc", "backspace":
				m.showFiles = false
				m.files = nil
			}
			return m, nil
		}

		// 詳細表示中の場合
		if m.showDetail {
			switch msg.String() {
			case "ctrl+c", "q":
				return m, tea.Quit
			case "esc", "backspace":
				m.showDetail = false
				m.detail = nil
			case "f":
				if m.detail != nil {
					files, err := ghclient.ListFiles(m.detail.Number)
					if err == nil {
						m.files = files
						m.showFiles = true
					}
				}
			}
			return m, nil
		}

		// 一覧表示中の場合
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
		case "enter":
			if len(m.prs) > 0 {
				detail, err := ghclient.GetPRDetail(m.prs[m.cursor].Number)
				if err == nil {
					m.detail = detail
					m.showDetail = true
				}
			}
		}
	}
	return m, nil
}

func (m model) View() string {
	// ファイル一覧表示モード
	if m.showFiles && m.files != nil {
		s := fmt.Sprintf("Files changed in PR #%d\n\n", m.detail.Number)
		for _, f := range m.files {
			s += fmt.Sprintf("  %s (+%d/-%d)\n", f.Path, f.Additions, f.Deletions)
		}
		s += "\nPress Esc to go back.\n"
		return s
	}

	// 詳細表示モード
	if m.showDetail && m.detail != nil {
		s := fmt.Sprintf("PR #%d: %s\n", m.detail.Number, m.detail.Title)
		s += fmt.Sprintf("State: %s\n", m.detail.State)
		s += fmt.Sprintf("Author: %s\n", m.detail.Author.Login)
		s += fmt.Sprintf("Created: %s\n", m.detail.CreatedAt)
		s += fmt.Sprintf("URL: %s\n", m.detail.URL)
		s += "\n--- Body ---\n"
		if m.detail.Body != "" {
			s += m.detail.Body
		} else {
			s += "(no description)"
		}
		s += "\n\nf: view files | Esc: go back\n"
		return s
	}

	// 一覧表示モード
	s := "Pull Requests\n\n"

	for i, pr := range m.prs {
		cursor := " "
		if m.cursor == i {
			cursor = ">"
		}
		s += fmt.Sprintf("%s #%d %s\n", cursor, pr.Number, pr.Title)
	}

	s += "\nEnter: view detail | q: quit\n"
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
