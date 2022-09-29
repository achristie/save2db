package listdbui

import (
	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

var docStyle = lipgloss.NewStyle().Margin(1, 2)

type item struct {
	title, desc string
}

func (i item) Title() string       { return i.title }
func (i item) Description() string { return i.desc }
func (i item) FilterValue() string { return i.title }

type model struct {
	list list.Model
}

func (m model) Init() tea.Cmd {
	return nil
}

type SelectMsg struct {
	item
}

func selectProjectCmd(s item) tea.Cmd {
	return func() tea.Msg {
		return SelectMsg{s}
	}
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		if msg.String() == "ctrl+c" || msg.String() == "esc" {
			return m, tea.Quit
		}
		if msg.String() == "enter" {
			// m.list.Items()[0].Title
			return m, selectProjectCmd((m.list.SelectedItem()).(item))
		}
	case tea.WindowSizeMsg:
		h, v := docStyle.GetFrameSize()
		m.list.SetSize(msg.Width-h, msg.Height-v)
	}

	var cmd tea.Cmd
	m.list, cmd = m.list.Update(msg)
	return m, cmd
}

func (m model) View() string {
	return docStyle.Render(m.list.View())
}

func New(items []string, title string) tea.Model {
	var l []list.Item
	for _, s := range items {
		l = append(l, item{title: s})
	}
	m := model{list: list.New(l, list.NewDefaultDelegate(), 0, 0)}
	m.list.Title = title
	m.list.SetFilteringEnabled(false)
	m.list.SetShowStatusBar(false)
	m.list.SetShowHelp(true)
	return m
}
