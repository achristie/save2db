package cmd

import (
	configure "github.com/achristie/save2db/pkg/tui/configure_db"
	"github.com/spf13/cobra"
)

func init() {
	configureCmd.AddCommand(configureDbCmd)
}

var configureDbCmd = &cobra.Command{
	Use:   "database",
	Short: "Configure your database",
	Run: func(cmd *cobra.Command, args []string) {
		configure.NewProgram()
	},
}
