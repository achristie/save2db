package cmd

import (
	"context"
	"fmt"
	"os"
	"time"

	"github.com/achristie/save2db/internal/tui/progress"
	"github.com/achristie/save2db/pkg/platts"
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

		// initialize Channel
		ch := make(chan platts.Result[platts.SymbolData])

		// setup TUI
		filters := make(map[string]string)
		filters["q"] = "= " + mdc
		filters["modifiedDate"] = ">= " + start
		main.p = tea.NewProgram(progress.New("FETCH SYMBOLS", filters))

		// fetch and store
		go func() {
			main.getSymbols(ctx, mdc, startDate, ch)
			writeToSvc(ctx, &main, ch, ss, false)
		}()

		// start TUI
		main.p.Start()
	},
}

func init() {
	fetchCmd.AddCommand(symCmd)
}
func (m *Main) getSymbols(ctx context.Context, mdc string, start time.Time, ch chan platts.Result[platts.SymbolData]) {
	m.client.GetReferenceData(start, 1000, mdc, ch)
}
