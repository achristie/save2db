package fetch

import (
	"github.com/charmbracelet/bubbles/progress"
	tea "github.com/charmbracelet/bubbletea"
)

func newProgress() progressModel {
	prog := progress.New(progress.WithScaledGradient("#FF7CCB", "#FDFF8C"))
	prog.Width = maxWidth

	return progressModel{progress: prog}
}

type progressModel struct {
	progress progress.Model
}

func (progressModel) Init() tea.Cmd {
	return nil
}

type ProgressMsg float64

func (m progressModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		return m, tea.Quit

	case ProgressMsg:
		cmd := m.progress.IncrPercent(float64(msg))
		m.progress.Update(msg)
		return m, cmd

	case progress.FrameMsg:
		progressModel, cmd := m.progress.Update(msg)
		m.progress = progressModel.(progress.Model)
		return m, cmd

	}
	return m, nil
}

func (m progressModel) View() string {
	// pad := strings.Repeat(" ", padding)
	return "\n" + m.progress.View() + "\n"
}
