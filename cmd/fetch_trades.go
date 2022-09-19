package cmd

import (
	"fmt"
	"log"
	"os"
	"strings"
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

		p := cli.NewProgram(fmt.Sprintf("Markets: %s, Modified Date >= [%s]", strings.Join(markets, ", "), start), []string{"Trades"})

		go func() {
			GetTrades(client, tds, markets, startDate, 1000, p)
		}()
		p.Start()
	},
}

func init() {
	fetchCmd.AddCommand(tradeCmd)
}

func GetTrades(client *platts.Client, db *TD.TradeDataStore, markets []string, start time.Time, pageSize int, p *tea.Program) {
	data := make(chan platts.Result[platts.TradeData])
	client.GetTradeData(markets, start, pageSize, data)
	p.Send(cli.StatusUpdater{Name: "Trades", Status: cli.Status{Category: cli.INPROGRESS, Msg: "In Progress"}})
	t := []platts.TradeResults{}

	for result := range data {
		if result.Err != nil {
			log.Printf("Error %s", result.Err)
			p.Send(cli.StatusUpdater{Name: "Trades", Status: cli.Status{Category: cli.ERROR, Msg: "An error occured, please retry."}})
			os.Exit(1)
		} else {
			res := result.Message
			pu := cli.ProgressUpdater{Name: "Trades", Percent: 1 / float64(res.Metadata.TotalPages)}
			p.Send(pu)
			log.Printf("Trade Data: %d records received from page [%d] in [%s] (%d total records). Adding to DB",
				len(res.Results), res.Metadata.Page, res.Metadata.QueryTime, res.Metadata.Count)
			t = append(t, res.Results...)
		}
	}
	if err := db.Add(t); err != nil {
		log.Printf("Error inserting records: %s", err)
		p.Send(cli.StatusUpdater{Name: "Trades", Status: cli.Status{Category: cli.ERROR, Msg: "An error occured, please retry."}})
	}
	p.Send(cli.StatusUpdater{Name: "Trades", Status: cli.Status{Category: cli.COMPLETED, Msg: fmt.Sprintf("Complete! Added [%d records] to [trades]", len(t))}})
}
