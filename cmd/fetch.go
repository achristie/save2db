package cmd

import (
	"fmt"
	"time"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var fetchCmd = &cobra.Command{
	Use:   "fetch [dataset]",
	Short: "Fetch data from Platts API.",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Printf("dataset not available, %s", args)
	},
}

func init() {

	fetchCmd.PersistentFlags().StringP("username", "u", "", "Your Username for calling Platts APIs")
	viper.BindPFlag("username", fetchCmd.PersistentFlags().Lookup("username"))

	fetchCmd.PersistentFlags().StringP("password", "p", "", "The Password associated with your Username")
	viper.BindPFlag("password", fetchCmd.PersistentFlags().Lookup("password"))

	fetchCmd.PersistentFlags().StringP("apikey", "a", "", "Your API Key for calling Platts APIs")
	viper.BindPFlag("apikey", fetchCmd.PersistentFlags().Lookup("apikey"))

	fetchCmd.PersistentFlags().StringP("mdc", "m", "ET", "Which Market Data Category to use")
	viper.BindPFlag("mdc", fetchCmd.PersistentFlags().Lookup("mdc"))

	fetchCmd.PersistentFlags().StringP("startDate", "t", time.Now().UTC().AddDate(0, 0, -7).Format("2006-01-02T15:04:05"), "Get updates since. modDate >= t")
	viper.BindPFlag("startDate", fetchCmd.PersistentFlags().Lookup("startDate"))

	rootCmd.AddCommand(fetchCmd)
}
