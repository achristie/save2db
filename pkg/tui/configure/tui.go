package configure

import (
	"github.com/achristie/save2db/pkg/tui/configure/input"
	"github.com/achristie/save2db/pkg/tui/configure/listdbui"
	tea "github.com/charmbracelet/bubbletea"
)

var p *tea.Program

type sessionState int

const (
	listView sessionState = iota
	inputView
)

// MainModel the main model of the program; holds other models and bubbles
type MainModel struct {
	state      sessionState
	list       tea.Model
	input      tea.Model
	windowSize tea.WindowSizeMsg
}

// View return the text UI to be output to the terminal
func (m MainModel) View() string {
	switch m.state {
	case listView:
		return m.list.View()
	default:
		return m.input.View()
	}
}

// New initialize the main model for your program
func New() MainModel {
	return MainModel{
		state: listView,
		list:  listdbui.New([]string{"SQLite", "PostgreSQL"}, "Select a Database"),
		input: input.NewConfigureDBModel(),
	}
}

// Init run any intial IO on program start
func (m MainModel) Init() tea.Cmd {
	return nil
}

// Update handle IO and commands
func (m MainModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd
	var cmds []tea.Cmd
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.windowSize = msg // pass this along to the entry view so it uses the full window size when it's initialized
	// case entryui.BackMsg:
	// 	m.state = projectView
	case listdbui.SelectMsg:
		m.selectedDB = msg.ActiveProjectID
		// m.input = entryui.New(m.er, m.activeProjectID, p, m.windowSize)
		m.state = inputView
	}

	switch m.state {
	case listView:
		newProject, newCmd := m.project.Update(msg)
		projectModel, ok := newProject.(projectui.Model)
		if !ok {
			panic("could not perform assertion on projectui model")
		}
		m.project = projectModel
		cmd = newCmd
	case inputView:
		newEntry, newCmd := m.input.Update(msg)
		entryModel, ok := newEntry.(input.Model)
		if !ok {
			panic("could not perform assertion on entryui model")
		}
		m.input = entryModel
		cmd = newCmd
	}
	cmds = append(cmds, cmd)
	return m, tea.Batch(cmds...)
}
