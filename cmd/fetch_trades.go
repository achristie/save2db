package cmd

import (
	"context"
	"fmt"
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

		// initialize TUI
		main.p = tui.NewProgram(fmt.Sprintf("Markets: %s, Modified Date >= [%s]", strings.Join(markets, ", "), start), []string{"Symbols"})

		// initialize Channel
		ch := make(chan platts.Result[platts.TradeData])

		go func() {
			main.getTrades(ctx, markets, startDate, ch)
			writeToSvc(ctx, &main, ch, ts)
		}()
		main.p.Start()
	},
}

func init() {
	fetchCmd.AddCommand(tradeCmd)
}

func (m *Main) getTrades(ctx context.Context, markets []string, start time.Time, ch chan platts.Result[platts.TradeData]) {
	m.client.GetTradeData(markets, start, 1000, ch)
	m.p.Send(tui.StatusUpdater{Name: "Symbols", Status: tui.Status{Category: tui.INPROGRESS, Msg: "In Progress"}})
}
