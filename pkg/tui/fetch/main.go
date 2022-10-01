package fetch

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type model struct {
	header   headerModel
	progress progressModel
	footer   footerModel
}

func New(title string, filters map[string]string) model {
	h := newHeader(title, filters)
	p := newProgress()
	f := newFooter()
	return model{header: h, progress: p, footer: f}
}

func (m model) Init() tea.Cmd {
	return nil
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg.(type) {
	case tea.KeyMsg:
		return m, tea.Quit
	}

	newProgress, pCmd := m.progress.Update(msg)
	m.progress = newProgress.(progressModel)

	newHeader, hCmd := m.header.Update(msg)
	m.header = newHeader.(headerModel)

	newFooter, fCmd := m.footer.Update(msg)
	m.footer = newFooter.(footerModel)

	return m, tea.Batch(pCmd, hCmd, fCmd)
}

func (m model) View() string {
	return lipgloss.JoinVertical(lipgloss.Center, "\n", m.header.View(), m.progress.View(), m.footer.View(), "\n")
}
