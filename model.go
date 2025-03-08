package main

import tea "github.com/charmbracelet/bubbletea"

type Model struct {
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
	return "Welcome to Arcane"
}
