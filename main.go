package main

import (
	"github.com/charmbracelet/bubbletea"
)

func main() {
	model := NewModel()
	program := tea.NewProgram(model, tea.WithAltScreen())
	if _, err := program.Run(); err != nil {
		panic(err)
	}
}
