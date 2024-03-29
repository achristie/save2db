package cmd

import (
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
		// initialize client
		main.client = platts.NewClient(config.Apikey, config.Username, config.Password)

		// initialize trade service
		ts, err := trades.New(db.GetDB(), config.Database.Name)
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
			main.getTrades(markets, startDate, ch)
			writeToSvc(&main, ch, ts, false)
		}()

		// start TUI
		main.p.Start()
	},
}

func init() {
	fetchCmd.AddCommand(tradeCmd)
}

func (m *Application) getTrades(markets []string, start time.Time, ch chan platts.Result[platts.TradeData]) {
	m.client.GetTradeData(markets, start, 1000, ch)
}
