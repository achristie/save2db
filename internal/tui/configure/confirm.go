package configure

import (
	tea "github.com/charmbracelet/bubbletea"
)

type confirmModel struct {
	selection string
	dict      map[string]string
}

func NewConfirmModel(selection string, dict map[string]string) tea.Model {
	m := confirmModel{selection: selection, dict: dict}
	return m
}

func (m confirmModel) Init() tea.Cmd {
	return nil
}

func (m confirmModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	return m, tea.Quit
}

func (m confirmModel) View() string {
	s := `
╭─────────────────────────────────────────╮
│                                         │
│    Database configuration is set        │
│                                         │
╰─────────────────────────────────────────╯`
	return "\n" + s + "\n"
}
