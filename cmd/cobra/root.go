package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var rootCmd = &cobra.Command{
	Use:   "platts",
	Short: "Platts CLI!",
	Long:  "Get your platts data via CLI",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Print("test")
	},
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func init() {
	cobra.OnInitialize(initConfig)
}

func initConfig() {

	home, err := os.UserHomeDir()
	cobra.CheckErr(err)
	viper.AddConfigPath(home)
	viper.SetConfigType("env")
	viper.SetConfigName(".platts")
	fmt.Println(home)

	viper.AutomaticEnv()
	if err := viper.ReadInConfig(); err == nil {
		fmt.Println("using config file: ", viper.ConfigFileUsed())
	}
}
