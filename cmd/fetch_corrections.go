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
		main.p = tui.NewProgram(fmt.Sprintf("MDC: [%s], Modified Date >= [%s]", mdc, start), []string{"Corrections"})

		// initialize assessments service
		as, err := services.NewAssessmentsService(ctx, db.GetDB(), config.DBSelection)
		if err != nil {
			fmt.Printf("assessments svc: %s", err)
			os.Exit(1)
		}
		main.assessmentService = as

		// initialize Channel
		main.chSymbolCorrection = make(chan platts.Result[platts.SymbolCorrection])

		go func() {
			main.getCorrections(ctx, startDate)
			main.writeCorrections(ctx)
		}()
		main.p.Start()
	},
}

func init() {
	fetchCmd.AddCommand(fcCmd)
}

func (m *Main) getCorrections(ctx context.Context, start time.Time) {
	m.client.GetDeletes(start, 10000, m.chSymbolCorrection)
	m.p.Send(tui.StatusUpdater{Name: "Corrections", Status: tui.Status{Category: tui.INPROGRESS, Msg: "In Progress"}})
}

func (m *Main) writeCorrections(ctx context.Context) {
	count := 0

	for result := range m.chSymbolCorrection {
		if result.Err != nil {
			log.Printf("fetch: %s", result.Err)
			m.p.Send(tui.StatusUpdater{Name: "Corrections", Status: tui.Status{Category: tui.ERROR, Msg: fmt.Sprint(result.Err)}})
			m.p.Quit()
		}

		res := result.Message
		m.p.Send(tui.ProgressUpdater{Name: "Corrections", Percent: 1 / float64(res.Metadata.TotalPages)})
		log.Printf("%+v", res.Metadata)
		for _, r := range res.Flatten() {
			_, err := m.assessmentService.Remove(ctx, m.tx, r)
			if err != nil {
				log.Printf("write: %s", err)
				m.p.Send(tui.StatusUpdater{Name: "Corrections", Status: tui.Status{Category: tui.ERROR, Msg: fmt.Sprint(err)}})
				m.p.Quit()
			}
			count += 1
		}
	}

	m.p.Send(tui.StatusUpdater{Name: "Corrections", Status: tui.Status{Category: tui.COMPLETED, Msg: fmt.Sprintf("Complete! Removed [%d records] from [assessments]", count)}})
	m.tx.Commit()
	m.p.Quit()
}
