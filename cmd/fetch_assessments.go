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

var faCmd = &cobra.Command{
	Use:   "assessments",
	Short: "Fetch assessment data",
	Long:  `Fetch assessments either by MDC (Market Data category) or Symbol(s) since t`,
	Run: func(cmd *cobra.Command, args []string) {
		ctx := context.Background()

		// initialize TUI
		main.p = tui.NewProgram(fmt.Sprintf("MDC: [%s], Modified Date >= [%s]", mdc, start), []string{"Assessments", "Deletes"})

		// initialize assessments service
		as, err := services.NewAssessmentsService(ctx, db.GetDB(), config.DBSelection)
		if err != nil {
			fmt.Printf("assessments svc: %s", err)
			os.Exit(1)
		}
		main.assessmentService = as

		// initialize Channel
		main.chSymbolHistory = make(chan platts.Result[platts.SymbolHistory])

		go func() {
			main.getAssessments(ctx, mdc, symbols, startDate)
			main.writeAssessments(ctx)
		}()
		main.p.Start()
	},
}

func init() {
	fetchCmd.AddCommand(faCmd)
}

func (m *Main) getAssessments(ctx context.Context, mdc string, symbols []string, start time.Time) {
	if len(symbols) > 0 {
		m.client.GetHistoryBySymbol(symbols, start, 10000, m.chSymbolHistory)
	} else {
		m.client.GetHistoryByMDC(mdc, start, 10000, m.chSymbolHistory)
	}
	m.p.Send(tui.StatusUpdater{Name: "Assessments", Status: tui.Status{Category: tui.INPROGRESS, Msg: "In Progress"}})
}

func (m *Main) writeAssessments(ctx context.Context) {
	count := 0

	for result := range m.chSymbolHistory {
		if result.Err != nil {
			log.Printf("fetch: %s", result.Err)
			m.p.Send(tui.StatusUpdater{Name: "Assessments", Status: tui.Status{Category: tui.ERROR, Msg: fmt.Sprint(result.Err)}})
			m.p.Quit()
		}

		res := result.Message
		m.p.Send(tui.ProgressUpdater{Name: "Assessments", Percent: 1 / float64(res.Metadata.TotalPages)})

		for _, r := range res.Flatten() {
			_, err := m.assessmentService.Add(ctx, m.tx, r)
			if err != nil {
				log.Printf("write: %s", err)
				m.p.Send(tui.StatusUpdater{Name: "Assessments", Status: tui.Status{Category: tui.ERROR, Msg: fmt.Sprint(err)}})
				m.p.Quit()
			}
			count += 1
		}
	}

	m.p.Send(tui.StatusUpdater{Name: "Assessments", Status: tui.Status{Category: tui.COMPLETED, Msg: fmt.Sprintf("Complete! Added [%d records] to [assessments]", count)}})
	m.tx.Commit()
	m.p.Quit()
}
