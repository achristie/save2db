package fetch

import (
	"fmt"

	tea "github.com/charmbracelet/bubbletea"
)

type headerModel struct {
	title   string
	filters map[string]string
}

func newHeader(title string, filters map[string]string) headerModel {
	return headerModel{title: title, filters: filters}
}

func (m headerModel) Init() tea.Cmd {
	return nil
}

func (m headerModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	return m, nil
}

func (m headerModel) View() string {
	return fmt.Sprintf("%s\n %+v", m.title, m.filters)
}
