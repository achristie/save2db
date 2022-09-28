package cmd

import (
	"fmt"

	tui "github.com/achristie/save2db/pkg/tui/configure"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/spf13/cobra"
)

func init() {
	configureCmd.AddCommand(configureDbCmd)
}

var configureDbCmd = &cobra.Command{
	Use:   "database",
	Short: "Configure your database",
	Run: func(cmd *cobra.Command, args []string) {
		// m := input.NewConfigureDBModel()
		// input.NewProgram(m)

		// m := listdbui.New([]string{"SQLite", "PostgreSQL", "SQL Server", "Oracle"}, "Select a Database")
		// if err := tea.NewProgram(m).Start(); err != nil {
		// 	fmt.Print(err)
		// }

		m := tui.New()
		if err := tea.NewProgram(m).Start(); err != nil {
			fmt.Print(err)
		}
	},
}
