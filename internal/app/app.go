package app

import (
	tea	"github.com/charmbracelet/bubbletea"

	"github.com/urabexon/prism/internal/ghclient"
	"github.com/urabexon/prism/internal/ui/prlist"
	"github.com/urabexon/prism/internal/ui/filelist"
)

type screen int

const (
	screenPRList screen = iota
	screenFileList
)

type Model struct {
	screen     screen
	prList     prlist.Model
	fileList   filelist.Model
	currentPR  int
}

func New(prs []ghclient.PR) Model {
	return Model{
		screen: screenPRList,
		prList: prlist.New(prs),
	}
}

func (m Model) Init() tea.Cmd {
	return nil
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch m.screen {
	case screenPRList:
		return m.updatePRList(msg)
	case screenFileList:
		return m.updateFileList(msg)
	}
	return m, nil
}

func (m Model) updatePRList(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		if msg.String() == "enter" {
			pr := m.prList.SelectedPR()
			files, err := ghclient.ListFiles(pr.Number)
			if err != nil {
				return m, nil // TODO: エラー処理
			}
			m.currentPR = pr.Number
			m.fileList = filelist.New(pr.Number, files)
			m.screen = screenFileList
			return m, nil
		}
	}
	updated, cmd := m.prList.Update(msg)
	m.prList = updated.(prlist.Model)
	return m, cmd
}

func (m Model) updateFileList(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		if msg.String() == "esc" {
			m.screen = screenPRList
			return m, nil
		}
	}
	updated, cmd := m.fileList.Update(msg)
	m.fileList = updated.(filelist.Model)
	return m, cmd
}

func (m Model) View() string {
	switch m.screen {
	case screenFileList:
		return m.fileList.View()
	default:
		return m.prList.View()
	}
}
