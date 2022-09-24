package cmd

import (
	"fmt"

	"github.com/achristie/save2db/internal/sqlite"
	"github.com/spf13/cobra"
	_ "modernc.org/sqlite"
)

var testCmd = &cobra.Command{
	Use:   "test",
	Short: "Fetch assessment data",
	Long: `Fetch assessments either by MDC (Market Data category) or Symbol(s) since t
	`,
	Run: func(cmd *cobra.Command, args []string) {
		db := sqlite.NewDB("database3.db")
		if err := db.Open(); err != nil {
			fmt.Print(err)
		}
	},
}

func init() {
	rootCmd.AddCommand(testCmd)
}
