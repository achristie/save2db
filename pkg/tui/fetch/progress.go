package fetch

import (
	"strings"

	"github.com/charmbracelet/bubbles/progress"
	tea "github.com/charmbracelet/bubbletea"
)

const (
	padding  = 2
	maxWidth = 80
)

func NewProgress() model {
	prog := progress.New(progress.WithScaledGradient("#FF7CCB", "#FDFF8C"))
	prog.Width = maxWidth

	return model{progress: prog}
}

type model struct {
	progress progress.Model
}

func (model) Init() tea.Cmd {
	return nil
}

type ProgressUpdater float64

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case ProgressUpdater:
		cmd := m.progress.IncrPercent(float64(msg))
		return m, cmd
	case progress.FrameMsg:
		progressModel, cmd := m.progress.Update(msg)
		m.progress = progressModel.(progress.Model)
		return m, cmd

	}
	return m, nil
}

func (m model) View() string {
	pad := strings.Repeat(" ", padding)
	return "\n" +
		pad + m.progress.View() + "\n"
}
