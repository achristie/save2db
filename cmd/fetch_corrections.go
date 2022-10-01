package cmd

import (
	"context"
	"fmt"
	"os"
	"time"

	"github.com/achristie/save2db/pkg/platts"
	tui "github.com/achristie/save2db/pkg/tui/progress"
	"github.com/achristie/save2db/services"
	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/spf13/cobra"
	_ "modernc.org/sqlite"
)

var fcCmd = &cobra.Command{
	Use:   "corrections",
	Short: "Fetch deleted assessment data",
	Long:  `Fetch corrections (deletes) since t`,
	Run: func(cmd *cobra.Command, args []string) {
		ctx := context.Background()

		// initialize TUI
		main.p = tui.NewProgram(fmt.Sprintf("MDC: [%s], Modified Date >= [%s]", mdc, start), []string{"Symbols"})

		// initialize assessments service
		as, err := services.NewAssessmentsService(ctx, db.GetDB(), config.DBSelection)
		if err != nil {
			fmt.Printf("assessments svc: %s", err)
			os.Exit(1)
		}

		// initialize Channel
		ch := make(chan platts.Result[platts.SymbolCorrection])

		go func() {
			main.getCorrections(ctx, startDate, ch)
			writeToSvc(ctx, &main, ch, as)
		}()
		main.p.Start()
	},
}

func init() {
	fetchCmd.AddCommand(fcCmd)
}

func (m *Main) getCorrections(ctx context.Context, start time.Time, ch chan platts.Result[platts.SymbolCorrection]) {
	m.client.GetDeletes(start, 10000, ch)
	m.p.Send(tui.StatusUpdater{Name: "Symbols", Status: tui.Status{Category: tui.INPROGRESS, Msg: "In Progress"}})
}
