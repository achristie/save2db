package cmd

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"os"

	"github.com/mitchellh/mapstructure"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var config Config

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

	// set defaults
	viper.SetDefault("Path", "database.db")
	viper.SetDefault("DBSelection", "SQLite")

	var result map[string]interface{}
	if err := viper.Unmarshal(&result); err != nil {
		fmt.Printf("decode map: %v", err)
	}

	if err := mapstructure.Decode(result, &config); err != nil {
		fmt.Printf("decode config: %v", err)
	}

}

type Database interface {
	Open() error
	BeginTx(context.Context, *sql.TxOptions) (*sql.Tx, error)
	GetDB() *sql.DB
}

type Config struct {
	username    string `mapstructure:"username"`
	apikey      string `mapstructure:"apikey"`
	password    string `mapstructure:"password"`
	dbHost      string `mapstructure:"dbhost"`
	dbPort      string `mapstructure:"dbport"`
	dbSelection string `mapstructure:"dbselection"`
	dbName      string `mapstructure:"dbname"`
	dbUsername  string `mapstructure:"dbusername"`
	dbPassword  string `mapstructure:"dbpassword"`
	path        string `mapstructure:"path"`
	fake        string `mapstructure:"fake"`
	errorLog    *log.Logger
	infoLog     *log.Logger
}
