package cmd

import (
	"context"
	"fmt"
	"os"
	"time"

	"github.com/achristie/save2db/pkg/platts"
	"github.com/achristie/save2db/pkg/tui/fetch"
	"github.com/achristie/save2db/services"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/spf13/cobra"
	_ "modernc.org/sqlite"
)

var symCmd = &cobra.Command{
	Use:   "symbols",
	Short: "Fetch symbol reference data",
	Run: func(cmd *cobra.Command, args []string) {
		ctx := context.Background()

		// initialize symbol service
		ss, err := services.NewSymbolService(ctx, db.GetDB(), config.DBSelection)
		if err != nil {
			fmt.Print(err)
			os.Exit(1)
		}

		filters := make(map[string]string)
		filters["yo"] = "oh now"
		main.p = tea.NewProgram(fetch.New("Fetch Symbols", filters))

		// initialize Channel
		ch := make(chan platts.Result[platts.SymbolData])

		go func() {
			main.getSymbols(ctx, mdc, startDate, ch)
			writeToSvc(ctx, &main, ch, ss)
		}()
		main.p.Start()
	},
}

func init() {
	fetchCmd.AddCommand(symCmd)
}
func (m *Main) getSymbols(ctx context.Context, mdc string, start time.Time, ch chan platts.Result[platts.SymbolData]) {
	m.client.GetReferenceData(start, 1000, mdc, ch)
	// m.p.Send(tui.StatusUpdater{Name: "Symbols", Status: tui.Status{Category: tui.INPROGRESS, Msg: "In Progress"}})
}
