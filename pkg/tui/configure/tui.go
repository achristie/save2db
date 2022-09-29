package tui

import (
	"github.com/achristie/save2db/pkg/tui/configure/input"
	"github.com/achristie/save2db/pkg/tui/configure/listdbui"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/spf13/viper"
)

type state int

const (
	listView state = iota
	inputView
)

// MainModel the main model of the program; holds other models and bubbles
type MainModel struct {
	state     state
	dblist    tea.Model
	input     tea.Model
	selection string
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
	switch msg := msg.(type) {
	case input.BackMsg:
		m.state = listView
	case input.SubmitMsg:
		for k, v := range msg.Values {
			viper.Set(k, v)
		}
		viper.WriteConfig()
		return m, tea.Quit
	case listdbui.SelectMsg:
		m.selection = msg.Title()
		if m.selection == "SQLite" {
			m.input = input.NewSQLiteModel()
		} else {
			m.input = input.NewConfigureDBModel()
		}
		m.state = inputView
	}

	if m.state == listView {
		m.dblist, cmd = m.dblist.Update(msg)
	} else {
		m.input, cmd = m.input.Update(msg)
	}
	return m, cmd
}
