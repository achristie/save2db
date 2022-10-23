package cmd

import (
	"bufio"
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

func init() {
	configureCmd.AddCommand(configureDb2Cmd)
}

type dbConfigModel struct {
	scanner *bufio.Scanner
}

var configureDb2Cmd = &cobra.Command{
	Use:   "database",
	Short: "configure database!",
	Run: func(cmd *cobra.Command, args []string) {
		m := dbConfigModel{
			scanner: bufio.NewScanner(os.Stdin),
		}
		var name, dsn string
		for {
			name = m.selectDb()
			dsn = m.getDSN(name)
			break
		}

		viper.Set("Database_Name", name)
		viper.Set("Database_DSN", dsn)
		err := viper.WriteConfig()
		if err != nil {
			cobra.CheckErr(err)
		}

		s := `
╭─────────────────────────────────────────╮
│    Database configuration is set        │
╰─────────────────────────────────────────╯`
		fmt.Print("\n" + s + "\n")
	},
}

func (m dbConfigModel) selectDb() string {
	fmt.Println("Select Database:")
	fmt.Println("(1) - SQLite")
	fmt.Println("(2) - PostgreSQL")

	m.scanner.Scan()
	switch m.scanner.Text() {
	case "1":
		return "SQLite"
	case "2":
		return "PostgreSQL"
	default:
		fmt.Printf("%q is not a supported selection\n", m.scanner.Text())
		return m.selectDb()
	}
}

func (m dbConfigModel) getDSN(name string) string {
	switch name {
	case "SQLite":
		fmt.Print("Enter SQLite Path: ")
		m.scanner.Scan()

		return m.scanner.Text()
	case "PostgreSQL":
		fmt.Print("Enter Host: ")
		m.scanner.Scan()
		host := m.scanner.Text()

		fmt.Print("Enter Port: ")
		m.scanner.Scan()
		port := m.scanner.Text()

		fmt.Print("Enter Username: ")
		m.scanner.Scan()
		un := m.scanner.Text()

		fmt.Print("Enter Password:")
		m.scanner.Scan()
		pw := m.scanner.Text()

		fmt.Print("Enter Database:")
		m.scanner.Scan()
		db := m.scanner.Text()

		return fmt.Sprintf("postgres://%s:%s@%s:%s/%s", un, pw, host, port, db)
	default:
		panic("unsupported")
	}

}
