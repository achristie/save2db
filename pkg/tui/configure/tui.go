package tui

import (
	"github.com/achristie/save2db/pkg/tui/configure/input"
	"github.com/achristie/save2db/pkg/tui/configure/listdbui"
	tea "github.com/charmbracelet/bubbletea"
)

type state int

const (
	listView state = iota
	inputView
)

// MainModel the main model of the program; holds other models and bubbles
type MainModel struct {
	state  state
	dblist tea.Model
	input  tea.Model
}

// View return the text UI to be output to the terminal
func (m MainModel) View() string {
	switch m.state {
	case listView:
		return m.dblist.View()
	default:
		return m.input.View()
	}
}

// New initialize the main model for your program
func New() MainModel {
	return MainModel{
		state:  listView,
		input:  input.NewConfigureDBModel(),
		dblist: listdbui.New([]string{"SQLite", "PostgreSQL"}, "Select a Database"),
	}
}

// Init run any intial IO on program start
func (m MainModel) Init() tea.Cmd {
	return nil
}

// Update handle IO and commands
func (m MainModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd
	switch msg.(type) {
	case input.BackMsg:
		m.state = listView
	case listdbui.SelectMsg:
		m.input = input.NewConfigureDBModel()
		m.state = inputView
	}

	if m.state == listView {
		m.dblist, cmd = m.dblist.Update(msg)
	} else {
		m.input, cmd = m.input.Update(msg)
	}
	return m, cmd
}
