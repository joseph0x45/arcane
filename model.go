package main

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/joseph0x45/arcane/models"
	"github.com/joseph0x45/arcane/shared"
)

const (
	authScreen = iota
	basicLoggedInScreen
)

type Model struct {
	screen     int
	authScreen *models.AuthScreenModel
}

func NewModel() *Model {
	return &Model{
		screen:     basicLoggedInScreen,
		authScreen: models.NewAuthScreenModel(),
	}
}

func (m *Model) Init() tea.Cmd {
	return nil
}

func (m *Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var (
		cmds []tea.Cmd
	)
	cmds = append(cmds, m.authScreen.Update(msg)...)
	switch msg.(type) {
	case tea.KeyMsg:
		keyMsgStr := msg.(tea.KeyMsg).String()
		switch keyMsgStr {
		case "ctrl+c":
			return m, tea.Quit
		case "ctrl+s":
			if m.screen == authScreen {
				m.screen = basicLoggedInScreen
				break
			}
			m.screen = authScreen
		}
		switch m.screen {
		case authScreen:
			m.authScreen.HandleKeyMessages(keyMsgStr)
		}
	}
	return m, tea.Batch(cmds...)
}

func (m *Model) View() string {
	switch m.screen {
	case authScreen:
		return m.authScreen.Render()
	case basicLoggedInScreen:
		return "This is the basic login screen"
	default:
		return shared.WelcomeASCIIArt + "\n\nSource code at https://github.com/joseph0x45/arcane"
	}
}
