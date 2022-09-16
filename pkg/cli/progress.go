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

func NewProgram(title string, names []string) *tea.Program {
	var pw []progressWrapper
	for _, s := range names {
		prog := progress.New(progress.WithScaledGradient("#FF7CCB", "#FDFF8C"))
		prog.Width = 60
		p := progressWrapper{name: s, progress: prog}
		pw = append(pw, p)
	}

	return tea.NewProgram(model{title: title, progress: pw})
}

var helpStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("6266262")).Render

// var errorStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("6266262")).Render
var ipStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("6266262")).Render
var successStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("6266262")).Render
var titleStyle = lipgloss.NewStyle().BorderStyle(lipgloss.NormalBorder()).BorderBottom(true).Foreground(lipgloss.Color("6266262")).Render

type ProgressUpdater struct {
	Name          string
	Percent       float64
	StatusMessage string
}

type progressWrapper struct {
	name      string
	percent   float64
	statusMsg string
	progress  progress.Model
}

type model struct {
	progress []progressWrapper
	title    string
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
				m.progress[i].statusMsg = msg.StatusMessage
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
	s += "\n" + pad + titleStyle(m.title)
	for _, v := range m.progress {
		s += "\n" + pad + v.name + " - " + statusString(v.statusMsg)
		s += "\n" + pad + m.progress[0].progress.ViewAs(v.percent) + "\n"
	}
	s += "\n\n" + pad + helpStyle("press any key to quit")
	return s

}

func statusString(msg string) string {
	switch msg {
	case "":
		return ipStyle("In Progress")
	default:
		return successStyle(msg)
	}
}
