package cmd

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/achristie/save2db/internal/sqlite"
	"github.com/achristie/save2db/pkg/cli"
	"github.com/achristie/save2db/pkg/platts"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	_ "modernc.org/sqlite"
)

type Main struct {
	client            *platts.Client
	tx                *sqlite.Tx
	p                 *tea.Program
	db                *sqlite.DB
	assessmentService *sqlite.AssessmentsService
	ch                chan platts.Result[platts.SymbolHistory]
}

var faCmd = &cobra.Command{
	Use:   "assessments",
	Short: "Fetch assessment data",
	Long:  `Fetch assessments either by MDC (Market Data category) or Symbol(s) since t`,
	Run: func(cmd *cobra.Command, args []string) {
		// create a platts api client
		ctx := context.Background()
		client := platts.NewClient(viper.GetString("apikey"), viper.GetString("username"), viper.GetString("password"))

		// initialize DB
		db := sqlite.NewDB("awc_database.db")
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
		as, err := sqlite.NewAssessmentsService(ctx, db)
		if err != nil {
			fmt.Print(err)
			os.Exit(1)
		}

		// initialize tui
		p := cli.NewProgram(fmt.Sprintf("MDC: [%s], Modified Date >= [%s]", mdc, start), []string{"Assessments", "Deletes"})

		// initialize Channel
		ch := make(chan platts.Result[platts.SymbolHistory])

		main := Main{client: client,
			db:                db,
			tx:                tx,
			p:                 p,
			assessmentService: as,
			ch:                ch,
		}

		go func() {
			main.getAssessments(ctx, mdc, symbols, startDate)
			main.writeAssessments(ctx)
			// getAssessments(ctx, client, tx, as, mdc, symbols, startDate, 10000, p)
			// getDeletes(client, as, startDate, 10000, p)
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
	m.p.Send(cli.StatusUpdater{Name: "Assessments", Status: cli.Status{Category: cli.INPROGRESS, Msg: "In Progress"}})
}

func (m *Main) writeAssessments(ctx context.Context) {
	su_error := cli.StatusUpdater{Name: "Assessments", Status: cli.Status{Category: cli.ERROR, Msg: "An error occured, please retry."}}
	count := 0

	for result := range m.ch {
		if result.Err != nil {
			log.Printf("Error! %s", result.Err)
			m.p.Send(su_error)
			m.p.Quit()
		}
		res := result.Message
		m.p.Send(cli.ProgressUpdater{Name: "Assessments", Percent: 1 / float64(res.Metadata.TotalPages)})
		// log.Printf("Assessment Data: %d records received from page [%d] in [%s] (%d total records).",
		// 	len(res.Results), res.Metadata.Page, res.Metadata.QueryTime, res.Metadata.Count)

		for _, r := range res.Flatten() {
			_, err := m.assessmentService.Add(ctx, m.tx, r)
			if err != nil {
				log.Printf("Error inserting records: %s", err)
				m.p.Send(su_error)
				m.p.Quit()
			}
			count += 1
		}
	}

	m.p.Send(cli.StatusUpdater{Name: "Assessments", Status: cli.Status{Category: cli.COMPLETED, Msg: fmt.Sprintf("Complete! Added [%d records] to [assessments]", count)}})
	m.tx.Commit()
	m.p.Quit()
}

// Get Deleted Assessments and remove from `assessments` table
// func getDeletes(client *platts.Client, db *MD.AssessmentsStore, start time.Time, pageSize int, p *tea.Program) {
// 	data := make(chan platts.Result[platts.SymbolCorrection])
// 	client.GetDeletes(start, pageSize, data)
// 	a := []platts.Assessment{}
// 	p.Send(cli.StatusUpdater{Name: "Deletes", Status: cli.Status{Category: cli.INPROGRESS, Msg: "In Progress"}})

// 	for result := range data {
// 		if result.Err != nil {
// 			log.Printf("Error! %s", result.Err)
// 		} else {
// 			res := result.Message
// 			p.Send(cli.ProgressUpdater{Name: "Deletes", Percent: 1 / float64(res.Metadata.TotalPages)})
// 			log.Printf("Deletes: %d records received from page [%d] in [%s] (%d total records).",
// 				len(res.Results), res.Metadata.Page, res.Metadata.QueryTime, res.Metadata.Count)

// 			a = append(a, res.Flatten()...)
// 		}
// 	}
// 	if err := db.Remove(a); err != nil {
// 		log.Printf("Error removing records: %s", err)
// 		p.Send(cli.StatusUpdater{Name: "Deletes", Status: cli.Status{Category: cli.ERROR, Msg: "An error occured, please retry."}})
// 	}
// 	p.Send(cli.StatusUpdater{Name: "Deletes", Status: cli.Status{Category: cli.COMPLETED, Msg: fmt.Sprintf("Complete! Removed [%d records] from [assessments]", len(a))}})
// }
