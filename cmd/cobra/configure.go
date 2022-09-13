package cmd

import (
	"bufio"
	"fmt"
	"os"

	"github.com/achristie/save2db/pkg/platts"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

func init() {
	rootCmd.AddCommand(configureCmd)
}

var configureCmd = &cobra.Command{
	Use:   "configure",
	Short: "Configure your Platts credentials.",
	Run: func(cmd *cobra.Command, args []string) {
		scanner := bufio.NewScanner(os.Stdin)
		var un, pw, apikey string
		for {
			fmt.Print("Enter Username: ")
			scanner.Scan()
			un = scanner.Text()

			fmt.Print("Enter Password: ")
			scanner.Scan()
			pw = scanner.Text()

			fmt.Print("Enter API Key: ")
			scanner.Scan()
			apikey = scanner.Text()

			break
		}

		fmt.Println("-----------------------")
		fmt.Println("Checking credentials...")
		fmt.Println("-----------------------")

		_, err := platts.GetToken(un, pw, apikey)
		if err != nil {
			fmt.Println(err)
			fmt.Println("Please try again.")
			os.Exit(1)
		}

		viper.Set("username", un)
		viper.Set("password", pw)
		viper.Set("apikey", apikey)
		viper.WriteConfig()
		fmt.Printf("Looks Good! Saved to config file [%s]", viper.ConfigFileUsed())

	},
}
