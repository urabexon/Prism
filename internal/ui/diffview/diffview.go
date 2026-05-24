package diffview

import (
    "fmt"
    "strings"

    tea "github.com/charmbracelet/bubbletea"
    "github.com/urabexon/prism/internal/ghclient"
)

type BackMsg struct{}

type Model struct {
	file   ghclient.FileDiff
	scroll int
    height int
}

func New(file ghclient.FileDiff) Model {
    return Model{file: file, height: 20}
}

func (m Model) Init() tea.Cmd { return nil }

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
    switch msg := msg.(type) {
    case tea.KeyMsg:
      switch msg.String() {
        case "j", "down":
            m.scroll++
        case "k", "up":
            if m.scroll > 0 {
                m.scroll--
            }
        case "esc", "q":
            return m, func() tea.Msg { return BackMsg{} }
        }
    }
    return m, nil
}

func (m Model) View() string {
    var b strings.Builder
    b.WriteString(fmt.Sprintf(" %s\n\n", m.file.FilePath()))

    for _, hunk := range m.file.Hunks {
        b.WriteString(hunk.Header + "\n")
        for _, line := range hunk.Lines {
            prefix := " "
            if line.Type == ghclient.LineAdded {
                prefix = "+"
            } else if line.Type == ghclient.LineRemoved {
                prefix = "-"
            }
            b.WriteString(prefix + line.Content + "\n")
        }
    }
    b.WriteString("\nj/k:scroll  esc:back")
    return b.String()
}
