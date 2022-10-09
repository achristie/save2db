package cmd

import (
	"context"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/achristie/save2db/internal/services/trades"
	"github.com/achristie/save2db/internal/tui/progress"
	"github.com/achristie/save2db/pkg/platts"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/spf13/cobra"
	_ "modernc.org/sqlite"
)

var tradeCmd = &cobra.Command{
	Use:   "trades",
	Short: "Fetch trade data (eWindow Market Data)",
	Run: func(cmd *cobra.Command, args []string) {
		ctx := context.Background()

		// initialize trade service
		ts, err := trades.New(ctx, db.GetDB(), config.DBSelection)
		if err != nil {
			fmt.Print(err)
			os.Exit(1)
		}

		// initialize Channel
		ch := make(chan platts.Result[platts.TradeData])

		// setup TUI
		filters := make(map[string]string)
		filters["markets"] = "in " + strings.Join(markets, ", ")
		filters["modifiedDate"] = ">= " + start
		main.p = tea.NewProgram(progress.New("FETCH SYMBOLS", filters))

		// fetch and store
		go func() {
			main.getTrades(ctx, markets, startDate, ch)
			writeToSvc(ctx, &main, ch, ts, false)
		}()

		// start TUI
		main.p.Start()
	},
}

func init() {
	fetchCmd.AddCommand(tradeCmd)
}

func (m *Main) getTrades(ctx context.Context, markets []string, start time.Time, ch chan platts.Result[platts.TradeData]) {
	m.client.GetTradeData(markets, start, 1000, ch)
}
