package progress

import (
	"fmt"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

var (
	statusBarStyle = lipgloss.NewStyle().
			Foreground(lipgloss.AdaptiveColor{Light: "#343433", Dark: "#C1C6B2"}).
			Background(lipgloss.AdaptiveColor{Light: "#D9DCCF", Dark: "#353533"})

	statusStyle = lipgloss.NewStyle().
			Inherit(statusBarStyle).
			Foreground(lipgloss.Color("#FFFDF5")).
			Background(lipgloss.Color("#FF5F87")).
			Padding(0, 1).
			MarginRight(1)

	statusText = lipgloss.NewStyle().Inherit(statusBarStyle)
)

type footerModel struct {
	statusMsg string
}

type StatusMsg struct {
	string
}

func StatusCmd(msg string) tea.Cmd {
	return func() tea.Msg {
		return StatusMsg{msg}
	}
}

func newFooter() footerModel {
	return footerModel{}
}

func (m footerModel) Init() tea.Cmd {
	return nil
}

func (m footerModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case StatusMsg:
		m.statusMsg = fmt.Sprint(msg)
		return m, nil
	}
	return m, nil
}

func (m footerModel) View() string {
	w := lipgloss.Width
	statusKey := statusStyle.Render("STATUS")
	statusVal := statusText.Copy().
		Width(maxWidth - w(statusKey)).
		Render(string(m.statusMsg))
	return lipgloss.JoinHorizontal(lipgloss.Left, statusKey, statusVal)
}
