package main

import tea "github.com/charmbracelet/bubbletea"

const (
	authScreen = iota
	basicLoggedInScreen
)

type Model struct {
	screen uint
}

func NewModel() *Model {
	return &Model{}
}

func (m *Model) Init() tea.Cmd {
	return nil
}

func (m *Model) Update(msg tea.Msg) (model tea.Model, cmd tea.Cmd) {
	switch msg.(type) {
	case tea.KeyMsg:
		switch msg.(tea.KeyMsg).String() {
		case "q":
			return m, tea.Quit
		}
	}
	return m, nil
}

func (m *Model) View() string {
	return arcaneASCIIArt + "\n\nhttps://github.com/joseph0x45/arcane"
}
