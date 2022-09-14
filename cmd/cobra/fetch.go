package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var fetchCmd = &cobra.Command{
	Use:   "fetch [dataset]",
	Short: "Fetch data from Platts API.",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("No dataset provided.")
	},
}

func init() {
	rootCmd.AddCommand(fetchCmd)

	fetchCmd.PersistentFlags().StringVar(&username, "username", "", "Your Username for calling Platts APIs")
	viper.BindPFlag("username", fetchCmd.PersistentFlags().Lookup("username"))

	fetchCmd.PersistentFlags().StringVar(&password, "password", "", "The Password associated with your Username")
	viper.BindPFlag("password", fetchCmd.PersistentFlags().Lookup("username"))

	fetchCmd.PersistentFlags().StringVar(&apikey, "apikey", "", "Your API Key for calling Platts APIs")
	viper.BindPFlag("apikey", fetchCmd.PersistentFlags().Lookup("username"))
}
