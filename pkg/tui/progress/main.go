package progress

import (
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

const (
	padding  = 2
	maxWidth = 60
)

var (
	subtle      = lipgloss.AdaptiveColor{Light: "#D9DCCF", Dark: "#383838"}
	borderStyle = lipgloss.NewStyle().
			Width(maxWidth).
			BorderStyle(lipgloss.NormalBorder()).
			BorderBottom(true).
			BorderForeground(subtle).
			MarginRight(2).
			Render
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
	pad := strings.Repeat(" ", padding)
	v1 := borderStyle(m.header.View())
	v2 := borderStyle(m.progress.View())
	v3 := m.footer.View()

	vj := lipgloss.JoinVertical(lipgloss.Left, "\n", v1, v2, v3, "\n")
	return lipgloss.JoinHorizontal(lipgloss.Left, pad, vj)
}
