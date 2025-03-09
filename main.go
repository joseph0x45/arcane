package main

import (
	"github.com/charmbracelet/bubbletea"
)

func main() {
	//check if user is logged in
	//if not ask them whether they
	//want to create an account or login
	model := NewModel()
	model.screen = authScreen
	program := tea.NewProgram(model, tea.WithAltScreen())
	if _, err := program.Run(); err != nil {
		panic(err)
	}
}
