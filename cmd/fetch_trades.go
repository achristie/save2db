package cmd

import (
	"log"
	"time"

	TD "github.com/achristie/save2db/internal/trade_data"
	"github.com/achristie/save2db/pkg/cli"
	"github.com/achristie/save2db/pkg/platts"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	_ "modernc.org/sqlite"
)

var tradeCmd = &cobra.Command{
	Use:   "trades",
	Short: "Fetch Trade Data (eWindow Market Data)",
	Run: func(cmd *cobra.Command, args []string) {

		// create a platts api client
		client := platts.NewClient(viper.GetString("apikey"), viper.GetString("username"), viper.GetString("password"))

		// initialize DB and create necessary tables
		db := TD.NewDb("database.db")
		tds := TD.NewTradeDataStore(db)

		p := cli.NewProgram("test", []string{"Trades"})

		go func() {
			GetTrades(client, tds, startDate, 1000, p)
		}()
		p.Start()
	},
}

func init() {
	fetchCmd.AddCommand(tradeCmd)
}

func GetTrades(client *platts.Client, db *TD.TradeDataStore, start time.Time, pageSize int, p *tea.Program) {
	data := make(chan platts.Result[platts.TradeData])
	client.GetTradeData(start, pageSize, data)

	for result := range data {
		if result.Err != nil {
			log.Printf("Error %s", result.Err)
		} else {
			res := result.Message
			pu := cli.ProgressUpdater{Name: "Trades", Percent: 1 / float64(res.Metadata.TotalPages)}
			p.Send(pu)
			log.Printf("Trade Data: %d records received from page [%d] in [%s] (%d total records). Adding to DB",
				len(res.Results), res.Metadata.Page, res.Metadata.QueryTime, res.Metadata.Count)
			if err := db.Add(res); err != nil {
				log.Printf("Error inserting records: %s", err)
			}
		}
	}
}
