package cli

import (
	"strings"

	"github.com/charmbracelet/bubbles/progress"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

const (
	padding  = 2
	maxWidth = 80
)

func NewProgram(names []string) *tea.Program {
	var pw []progressWrapper
	for _, s := range names {
		prog := progress.New(progress.WithScaledGradient("#FF7CCB", "#FDFF8C"))
		prog.Width = 80
		p := progressWrapper{name: s, progress: prog}
		pw = append(pw, p)
	}

	return tea.NewProgram(model{progress: pw})
}

var helpStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("6266262")).Render

type ProgressUpdater struct {
	Name    string
	Percent float64
}

type progressWrapper struct {
	name     string
	percent  float64
	progress progress.Model
}

type model struct {
	progress []progressWrapper
}

func (model) Init() tea.Cmd {
	return nil
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		return m, tea.Quit

	case ProgressUpdater:
		for i := range m.progress {
			if m.progress[i].name == msg.Name {
				m.progress[i].percent += msg.Percent
				break
			}
		}
		return m, nil

	default:
		return m, nil
	}
}

func (m model) View() string {
	pad := strings.Repeat(" ", padding)
	s := ""
	for _, v := range m.progress {
		s += "\n" + pad + v.name + pad + m.progress[0].progress.ViewAs(v.percent) + "\n"
	}
	s += "\n\n" + pad + helpStyle("press any key to quit")
	return s

}
