package progress

import (
	"fmt"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

var (
	titleStyle = lipgloss.NewStyle().
			MarginRight(5).
			Padding(0, 1).
			Foreground(lipgloss.Color("#FFFDF5")).
			Background(lipgloss.Color("#6495ED")).
			Render

	filterStyle = lipgloss.NewStyle().MarginLeft(2).Render
)

type headerModel struct {
	title   string
	filters map[string]string
	f       string
}

func newHeader(title string, filters map[string]string) headerModel {
	var f string
	for k, v := range filters {
		f += fmt.Sprintf("\n%s: [%s]", k, v)
	}

	return headerModel{title: title, filters: filters, f: f}
}

func (m headerModel) Init() tea.Cmd {
	return nil
}

func (m headerModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	return m, nil
}

func (m headerModel) View() string {
	s := fmt.Sprintf("%s%s", titleStyle(m.title), filterStyle(m.f))
	return s
}
