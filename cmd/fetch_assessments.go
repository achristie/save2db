package cmd

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/achristie/save2db/pg"
	"github.com/achristie/save2db/pkg/platts"
	tui "github.com/achristie/save2db/pkg/tui/progress"
	"github.com/achristie/save2db/services"
	"github.com/achristie/save2db/sqlite"
	tea "github.com/charmbracelet/bubbletea"
	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/spf13/cobra"
	_ "modernc.org/sqlite"
)

type Main struct {
	client            *platts.Client
	tx                *sql.Tx
	p                 *tea.Program
	assessmentService *services.AssessmentsService
	ch                chan platts.Result[platts.SymbolHistory]
}

var faCmd = &cobra.Command{
	Use:   "assessments",
	Short: "Fetch assessment data",
	Long:  `Fetch assessments either by MDC (Market Data category) or Symbol(s) since t`,
	Run: func(cmd *cobra.Command, args []string) {
		// create a platts api client
		ctx := context.Background()
		client := platts.NewClient(config.APIKey, config.Username, config.Password)

		var db Database
		switch config.DBSelection {
		case "PostgreSQL":
			conn := fmt.Sprintf("postgres://%s:%s@%s:%s/%s", config.DBUsername, config.DBPassword,
				config.DBHost, config.DBPort, config.DBName)
			db = pg.NewDB(conn)
		default:
			db = sqlite.NewDB(config.Path)
		}

		if err := db.Open(); err != nil {
			fmt.Print(err)
			os.Exit(1)
		}

		// begin a transaction
		tx, err := db.BeginTx(ctx, nil)
		if err != nil {
			fmt.Print(err)
		}

		// initialize assessments service
		as, err := services.NewAssessmentsService(ctx, db.GetDB(), config.DBSelection)
		if err != nil {
			fmt.Print(err)
			os.Exit(1)
		}

		// initialize TUI
		p := tui.NewProgram(fmt.Sprintf("MDC: [%s], Modified Date >= [%s]", mdc, start), []string{"Assessments", "Deletes"})

		// initialize Channel
		ch := make(chan platts.Result[platts.SymbolHistory])

		main := Main{client: client,
			tx:                tx,
			p:                 p,
			assessmentService: as,
			ch:                ch,
		}

		go func() {
			main.getAssessments(ctx, mdc, symbols, startDate)
			main.writeAssessments(ctx)
		}()
		p.Start()
	},
}

func init() {
	fetchCmd.AddCommand(faCmd)
}

func (m *Main) getAssessments(ctx context.Context, mdc string, symbols []string, start time.Time) {
	if len(symbols) > 0 {
		m.client.GetHistoryBySymbol(symbols, start, 10000, m.ch)
	} else {
		m.client.GetHistoryByMDC(mdc, start, 10000, m.ch)
	}
	m.p.Send(tui.StatusUpdater{Name: "Assessments", Status: tui.Status{Category: tui.INPROGRESS, Msg: "In Progress"}})
}

func (m *Main) writeAssessments(ctx context.Context) {
	count := 0

	for result := range m.ch {
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

func (m *Main) getDeletes(ctx context.Context, start time.Time) {
	// if len(symbols) > 0 {
	m.client.GetDeletes(mdc, start, 10000, m.ch)
	// } else {
	// 	m.client.GetHistoryByMDC(mdc, start, 10000, m.ch)
	// }
	m.p.Send(tui.StatusUpdater{Name: "Deletes", Status: tui.Status{Category: tui.INPROGRESS, Msg: "In Progress"}})
}

// Get Deleted Assessments and remove from `assessments` table
func getDeletes() {
	data := make(chan platts.Result[platts.SymbolCorrection])
	client.GetDeletes(start, pageSize, data)
	a := []platts.Assessment{}
	p.Send(cli.StatusUpdater{Name: "Deletes", Status: cli.Status{Category: cli.INPROGRESS, Msg: "In Progress"}})

	for result := range data {
		if result.Err != nil {
			log.Printf("Error! %s", result.Err)
		} else {
			res := result.Message
			p.Send(cli.ProgressUpdater{Name: "Deletes", Percent: 1 / float64(res.Metadata.TotalPages)})
			log.Printf("Deletes: %d records received from page [%d] in [%s] (%d total records).",
				len(res.Results), res.Metadata.Page, res.Metadata.QueryTime, res.Metadata.Count)

			a = append(a, res.Flatten()...)
		}
	}
	if err := db.Remove(a); err != nil {
		log.Printf("Error removing records: %s", err)
		p.Send(cli.StatusUpdater{Name: "Deletes", Status: cli.Status{Category: cli.ERROR, Msg: "An error occured, please retry."}})
	}
	p.Send(cli.StatusUpdater{Name: "Deletes", Status: cli.Status{Category: cli.COMPLETED, Msg: fmt.Sprintf("Complete! Removed [%d records] from [assessments]", len(a))}})
}
