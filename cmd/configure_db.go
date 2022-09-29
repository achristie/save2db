package cmd

import (
	"fmt"

	"github.com/achristie/save2db/pkg/tui/configure"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

func init() {
	configureCmd.AddCommand(configureDbCmd)
}

var configureDbCmd = &cobra.Command{
	Use:   "database",
	Short: "Configure your database",
	Run: func(cmd *cobra.Command, args []string) {
		m := NewModel()
		if err := tea.NewProgram(m).Start(); err != nil {
			fmt.Print(err)
		}

	},
}

type state int

const (
	listView state = iota
	inputView
	confirmView
)

type model struct {
	state     state
	list      tea.Model
	input     tea.Model
	confirm   tea.Model
	selection string
}

func NewModel() model {
	return model{
		state: listView,
		list:  configure.NewList([]string{"SQLite", "PostgreSQL"}, "Select a Database"),
	}
}

func (m model) Init() tea.Cmd {
	return nil
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd
	switch msg := msg.(type) {
	case configure.BackMsg:
		m.state = listView
	case configure.SubmitMsg:
		// store values in config file
		for k, v := range msg.Values {
			viper.Set(k, v)
		}
		viper.Set("DBSelection", m.selection)
		viper.WriteConfig()

		m.state = confirmView
		m.confirm = configure.NewConfirmModel(m.selection, msg.Values)
		return m, nil
	case configure.SelectMsg:
		m.selection = msg.Title()

		if m.selection == "SQLite" {
			m.input = configure.NewSQLiteModel()
		} else {
			m.input = configure.NewConfigureDBModel()
		}
		m.state = inputView
	}

	switch m.state {
	case listView:
		m.list, cmd = m.list.Update(msg)
	case inputView:
		m.input, cmd = m.input.Update(msg)
	case confirmView:
		m.confirm, cmd = m.confirm.Update(msg)
	}
	return m, cmd
}

func (m model) View() string {
	switch m.state {
	case listView:
		return m.list.View()
	case inputView:
		return m.input.View()
	case confirmView:
		return m.confirm.View()
	default:
		return ""
	}
}
