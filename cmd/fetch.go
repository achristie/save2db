package cmd

import (
	"fmt"
	"time"

	"github.com/achristie/save2db/internal/token"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var (
	main      Application
	db        Database
	start     string
	startDate time.Time
	mdc       string
	symbols   []string
	markets   []string
	watchlist string
)

var fetchCmd = &cobra.Command{
	Use:   "fetch [dataset]",
	Short: "Fetch data from Platts API.",
	PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
		// execute root
		rootCmd.PersistentPreRun(cmd, args)

		// parse the date flag
		var err error
		startDate, err = ParseDate(start)
		if err != nil {
			return err
		}

		// fetch requires a token anyway so lets get one now
		if err := InitToken(config); err != nil {
			return err
		}

		return nil
		// return InitDB(config)
	},
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Printf("dataset not available, %s", args)
	},
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
	fetchCmd.PersistentFlags().StringVarP(&watchlist, "watchlist", "w", "", "Use Symbols from a PDP Watchlist.")
	fetchCmd.MarkFlagsMutuallyExclusive("mdc", "symbol", "watchlist")

	fetchCmd.PersistentFlags().StringVarP(&start, "startDate", "t", time.Now().UTC().AddDate(0, 0, -7).Format("2006-01-02T15:04:05"), "Get assessments since t. Ex. 2021-01-01T00:00:00")

	fetchCmd.PersistentFlags().StringSliceVarP(&markets, "markets", "m", nil, "Markets to get Trades for. Ex: 'EU BFOE, US Midwest'")

	rootCmd.AddCommand(fetchCmd)
}

func ParseDate(start string) (time.Time, error) {
	startDate, err := time.Parse("2006-01-02T15:04:05", start)
	if err != nil {
		return time.Time{}, err
	}

	return startDate, nil
}

func InitToken(cfg Config) error {
	tc := token.NewTokenClient(cfg.Username, cfg.Password, cfg.Apikey)
	_, err := tc.GetToken()
	if err != nil {
		return err
	}

	return nil
}
