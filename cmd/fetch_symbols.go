package cmd

import (
	"context"
	"fmt"
	"os"
	"time"

	"github.com/achristie/save2db/pkg/platts"
	tui "github.com/achristie/save2db/pkg/tui/progress"
	"github.com/achristie/save2db/services"
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

		// initialize TUI
		main.p = tui.NewProgram(fmt.Sprintf("MDC: [%s], Modified Date >= [%s]", mdc, start), []string{"Symbols"})

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
	m.p.Send(tui.StatusUpdater{Name: "Symbols", Status: tui.Status{Category: tui.INPROGRESS, Msg: "In Progress"}})
}
