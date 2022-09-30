package cmd

import (
	"context"
	"fmt"
	"log"
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
		main.symbolService = ss

		// initialize TUI
		main.p = tui.NewProgram(fmt.Sprintf("MDC: [%s], Modified Date >= [%s]", mdc, start), []string{"Symbols"})

		// initialize Channel
		main.chSymbolData = make(chan platts.Result[platts.SymbolData])

		go func() {
			main.getSymbols(ctx, mdc, startDate)
			main.writeSymbols(ctx)
		}()
		main.p.Start()
	},
}

func init() {
	fetchCmd.AddCommand(symCmd)
}
func (m *Main) getSymbols(ctx context.Context, mdc string, start time.Time) {
	m.client.GetReferenceData(start, 1000, mdc, m.chSymbolData)
	m.p.Send(tui.StatusUpdater{Name: "Symbols", Status: tui.Status{Category: tui.INPROGRESS, Msg: "In Progress"}})
}

func (m *Main) writeSymbols(ctx context.Context) {
	count := 0

	for result := range m.chSymbolData {
		if result.Err != nil {
			log.Printf("fetch: %s", result.Err)
			m.p.Send(tui.StatusUpdater{Name: "Symbols", Status: tui.Status{Category: tui.ERROR, Msg: fmt.Sprint(result.Err)}})
			m.p.Quit()
		}

		res := result.Message
		m.p.Send(tui.ProgressUpdater{Name: "Symbols", Percent: 1 / float64(res.Metadata.TotalPages)})

		for _, r := range res.Results {
			_, err := m.symbolService.Add(ctx, m.tx, r)
			if err != nil {
				log.Printf("write: %s", err)
				m.p.Send(tui.StatusUpdater{Name: "Symbols", Status: tui.Status{Category: tui.ERROR, Msg: fmt.Sprint(err)}})
				m.p.Quit()
			}
			count += 1
		}
	}
	m.p.Send(tui.StatusUpdater{Name: "Symbols", Status: tui.Status{Category: tui.COMPLETED, Msg: fmt.Sprintf("Complete! Added [%d records] to [assessments]", count)}})
	m.tx.Commit()
	m.p.Quit()

}
