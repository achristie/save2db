package cmd

import (
	"fmt"

	"github.com/achristie/save2db/pkg/platts"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var viewCmd = &cobra.Command{
	Use:   "view",
	Short: "view...",
	Run: func(cmd *cobra.Command, args []string) {
		// fmt.Printf("dataset not available, %s", args)
		client := platts.NewClient(config.APIKey, config.Username, config.Password)

		wl, err := client.GetWatchlistByName("Watchlist")
		if err != nil {
			fmt.Print(err)
		}
		fmt.Printf("%+v", wl)
	},
}

func init() {

	viewCmd.PersistentFlags().StringP("username", "u", "", "Your username for calling Platts APIs")
	viper.BindPFlag("username", viewCmd.PersistentFlags().Lookup("username"))

	viewCmd.PersistentFlags().StringP("password", "p", "", "The password associated with your Username")
	viper.BindPFlag("password", viewCmd.PersistentFlags().Lookup("password"))

	viewCmd.PersistentFlags().StringP("apikey", "a", "", "Your API Key for calling Platts APIs")
	viper.BindPFlag("apikey", viewCmd.PersistentFlags().Lookup("apikey"))

	rootCmd.AddCommand(viewCmd)
}
