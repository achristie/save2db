package cmd

import (
	"context"
	"fmt"
	"log"
	"os"
	"strings"
	"time"

	"github.com/achristie/save2db/pkg/platts"
	tui "github.com/achristie/save2db/pkg/tui/progress"
	"github.com/achristie/save2db/services"
	"github.com/spf13/cobra"
	_ "modernc.org/sqlite"
)

var tradeCmd = &cobra.Command{
	Use:   "trades",
	Short: "Fetch trade data (eWindow Market Data)",
	Run: func(cmd *cobra.Command, args []string) {
		ctx := context.Background()

		// initialize symbol service
		ts, err := services.NewTradeService(ctx, db.GetDB(), config.DBSelection)
		if err != nil {
			fmt.Print(err)
			os.Exit(1)
		}
		main.tradeService = ts

		// initialize TUI
		main.p = tui.NewProgram(fmt.Sprintf("Markets: %s, Modified Date >= [%s]", strings.Join(markets, ", "), start), []string{"Trades"})

		// initialize Channel
		main.chTradeData = make(chan platts.Result[platts.TradeData])

		go func() {
			main.getTrades(ctx, markets, startDate)
			main.writeTrades(ctx)
		}()
		main.p.Start()
	},
}

func init() {
	fetchCmd.AddCommand(tradeCmd)
}

func (m *Main) getTrades(ctx context.Context, markets []string, start time.Time) {
	m.client.GetTradeData(markets, start, 1000, m.chTradeData)
	m.p.Send(tui.StatusUpdater{Name: "Trades", Status: tui.Status{Category: tui.INPROGRESS, Msg: "In Progress"}})
}

func (m *Main) writeTrades(ctx context.Context) {
	count := 0

	for result := range m.chTradeData {
		if result.Err != nil {
			log.Printf("fetch: %s", result.Err)
			m.p.Send(tui.StatusUpdater{Name: "Trades", Status: tui.Status{Category: tui.ERROR, Msg: fmt.Sprint(result.Err)}})
			m.p.Quit()
		}

		res := result.Message
		m.p.Send(tui.ProgressUpdater{Name: "Trades", Percent: 1 / float64(res.Metadata.TotalPages)})

		for _, r := range res.Results {
			_, err := m.tradeService.Add(ctx, m.tx, r)
			if err != nil {
				log.Printf("write: %s", err)
				m.p.Send(tui.StatusUpdater{Name: "Trades", Status: tui.Status{Category: tui.ERROR, Msg: fmt.Sprint(err)}})
				m.p.Quit()
			}
			count += 1
		}
	}
	m.p.Send(tui.StatusUpdater{Name: "Trades", Status: tui.Status{Category: tui.COMPLETED, Msg: fmt.Sprintf("Complete! Added [%d records] to [trades]", count)}})
	m.tx.Commit()
	m.p.Quit()
}
