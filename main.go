package main

import (
	"github.com/charmbracelet/bubbletea"
)

func main() {
	//check if user is logged in
	//if not display authScreen
	model := NewModel()
	program := tea.NewProgram(model, tea.WithAltScreen())
	if _, err := program.Run(); err != nil {
		panic(err)
	}
}
