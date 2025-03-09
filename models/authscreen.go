package models

import (
	"net/mail"

	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

var (
	appNameStyle     = lipgloss.NewStyle().Background(lipgloss.Color("20")).Padding(0, 1)
	errorTextStyle   = lipgloss.NewStyle().Foreground(lipgloss.Color("9"))
	screenTitleStyle = lipgloss.NewStyle().Bold(true)
)

type AuthScreenModel struct {
	EmailInput    textinput.Model
	PasswordInput textinput.Model
	errorStr      string
}

func NewAuthScreenModel() *AuthScreenModel {
	emailInput := textinput.New()
	emailInput.Placeholder = "Enter your email"
	emailInput.Focus()
	passwordInput := textinput.New()
	passwordInput.Placeholder = "Enter your password"
	return &AuthScreenModel{
		EmailInput:    emailInput,
		PasswordInput: passwordInput,
	}
}

func (m *AuthScreenModel) HandleKeyMessages(keyMsg string) {
	if keyMsg == "enter" {
		_, err := mail.ParseAddress(m.EmailInput.Value())
		if err != nil {
			m.errorStr = "Invalid Email"
		}
		if m.PasswordInput.Value() == "" || len(m.PasswordInput.Value()) > 72 {
			m.errorStr += "\nPassword can't be empty or more than 72 characters"
		}
	}
	if keyMsg == "up" || keyMsg == "down" {
		m.SwitchFocus()
	}
}

func (m *AuthScreenModel) Update(msg tea.Msg) []tea.Cmd {
	var (
		cmds []tea.Cmd
		cmd  tea.Cmd
	)
	m.EmailInput, cmd = m.EmailInput.Update(msg)
	cmds = append(cmds, cmd)
	m.PasswordInput, cmd = m.PasswordInput.Update(msg)
	cmds = append(cmds, cmd)
	return cmds
}

func (m *AuthScreenModel) SwitchFocus() {
	if m.EmailInput.Focused() {
		m.EmailInput.Blur()
		m.PasswordInput.Focus()
	} else {
		m.EmailInput.Focus()
		m.PasswordInput.Blur()
	}
}

func (m *AuthScreenModel) Render() string {
	output := appNameStyle.Render("Authenticate or create an account to continue.") + "\n\n"
	output += m.EmailInput.View()
	output += "\n\n"
	output += m.PasswordInput.View()
	if m.errorStr != "" {
		output += "\n\n" + errorTextStyle.Render(m.errorStr)
	}
	output += "\n\n"
	output += "Use up/down to switch focus  |  <Ctrl+Enter> to submit"
	return output
}
