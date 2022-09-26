package cmd

import (
	"context"
	"database/sql"
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var rootCmd = &cobra.Command{
	Use:              "platts-cli",
	Short:            "Platts CLI!",
	Long:             "Get your platts data via CLI",
	PersistentPreRun: initLogging,
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
	configName := ".plattsrc"
	viper.AddConfigPath(home)
	viper.SetConfigType("env")
	viper.SetConfigName(configName)

	// create if the config does not yet exists
	os.OpenFile(fmt.Sprintf("%s/%s", home, configName), os.O_CREATE|os.O_RDONLY, 0666)

	// problematic to use env variables
	// viper.AutomaticEnv()

	err = viper.ReadInConfig()
	if err != nil {
		fmt.Printf("unable to read config: %v", err)
	}

}

type Database interface {
	Open() error
	BeginTx(context.Context, *sql.TxOptions) (*sql.Tx, error)
	GetDB() *sql.DB
}
