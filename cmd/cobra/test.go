package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

func init() {
	rootCmd.AddCommand(testCmd)
}

var testCmd = &cobra.Command{
	Use:   "test",
	Short: "Configure your Platts credentials.",
	Run: func(cmd *cobra.Command, args []string) {
		un := viper.Get("username")
		fmt.Print(un)
	},
}
