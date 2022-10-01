package cmd

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/achristie/save2db/pg"
	"github.com/achristie/save2db/pkg/platts"
	"github.com/achristie/save2db/services"
	"github.com/achristie/save2db/sqlite"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

type Main struct {
	client             *platts.Client
	tx                 *sql.Tx
	p                  *tea.Program
	tradeService       *services.TradeService
	chTradeData        chan platts.Result[platts.TradeData]
	chSymbolCorrection chan platts.Result[platts.SymbolCorrection]
}

var (
	main      Main
	db        Database
	start     string
	startDate time.Time
	mdc       string
	symbols   []string
	markets   []string
)

var fetchCmd = &cobra.Command{
	Use:   "fetch [dataset]",
	Short: "Fetch data from Platts API.",
	PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
		// execute root
		rootCmd.PersistentPreRun(cmd, args)

		// parse the date flag
		var err error
		startDate, err = time.Parse("2006-01-02T15:04:05", start)
		if err != nil {
			return err
		}

		// fetch requires a token anyway so lets get one now
		_, err = platts.GetToken(viper.GetString("username"), viper.GetString("password"), viper.GetString("apikey"))
		if err != nil {
			return fmt.Errorf("invalid credentials. Did you use the `configure` command?")
		}

		return InitDB()
	},
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Printf("dataset not available, %s", args)
	},
}

// setup db stuff
func InitDB() error {
	ctx := context.Background()
	client := platts.NewClient(config.APIKey, config.Username, config.Password)

	switch config.DBSelection {
	case "PostgreSQL":
		conn := fmt.Sprintf("postgres://%s:%s@%s:%s/%s", config.DBUsername, config.DBPassword,
			config.DBHost, config.DBPort, config.DBName)
		db = pg.NewDB(conn)
	default:
		db = sqlite.NewDB(config.Path)
	}

	if err := db.Open(); err != nil {
		return fmt.Errorf("db open: %w", err)
	}

	// begin a transaction
	tx, err := db.BeginTx(ctx, nil)
	if err != nil {
		return fmt.Errorf("db tx: %w", err)
	}

	main = Main{client: client,
		tx: tx,
	}
	return nil
}

func init() {

	fetchCmd.PersistentFlags().StringP("username", "u", "", "Your username for calling Platts APIs")
	viper.BindPFlag("username", fetchCmd.PersistentFlags().Lookup("username"))

	fetchCmd.PersistentFlags().StringP("password", "p", "", "The password associated with your Username")
	viper.BindPFlag("password", fetchCmd.PersistentFlags().Lookup("password"))

	fetchCmd.PersistentFlags().StringP("apikey", "a", "", "Your API Key for calling Platts APIs")
	viper.BindPFlag("apikey", fetchCmd.PersistentFlags().Lookup("apikey"))

	fetchCmd.PersistentFlags().StringVar(&mdc, "mdc", "", "Market Data Category to get assessments for. Ex: IO")
	fetchCmd.PersistentFlags().StringSliceVarP(&symbols, "symbol", "s", nil, "Symbols to get assessments for. Ex: 'PCAAS00, PCAAT00'")
	fetchCmd.MarkFlagsMutuallyExclusive("mdc", "symbol")

	fetchCmd.PersistentFlags().StringVarP(&start, "startDate", "t", time.Now().UTC().AddDate(0, 0, -7).Format("2006-01-02T15:04:05"), "Get assessments since t. Ex. 2021-01-01T00:00:00")

	fetchCmd.PersistentFlags().StringSliceVarP(&markets, "markets", "m", nil, "Markets to get Trades for. Ex: 'EU BFOE, US Midwest'")

	rootCmd.AddCommand(fetchCmd)
}
