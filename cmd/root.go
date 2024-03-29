package cmd

import (
	"database/sql"
	"fmt"
	"log"
	"os"

	"github.com/achristie/save2db/pkg/platts"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/mitchellh/mapstructure"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var config Config

type Config struct {
	Username string `mapstructure:"username"`
	Apikey   string `mapstructure:"apikey"`
	Password string `mapstructure:"password"`
	Fake     string `mapstructure:"fake"`
	Database DB     `mapstructure:",squash"`
}

type DB struct {
	Name string `mapstructure:"database_name"`
	DSN  string `mapstructure:"database_dsn"`
}

type Application struct {
	client *platts.Client
	tx     *sql.Tx
	p      *tea.Program
	logger *log.Logger
}

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

	err = viper.ReadInConfig()
	if err != nil {
		fmt.Printf("unable to read config: %v", err)
	}

	// set defaults
	viper.SetDefault("Database_Name", "SQLite")
	viper.SetDefault("Database_DSN", "database.db")

	var result map[string]interface{}
	if err := viper.Unmarshal(&result); err != nil {
		fmt.Printf("decode map: %v", err)
	}

	if err := mapstructure.Decode(result, &config); err != nil {
		fmt.Printf("decode config: %v", err)
	}
}
