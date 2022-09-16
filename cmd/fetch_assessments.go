package cmd

import (
	"fmt"
	"log"
	"time"

	MD "github.com/achristie/save2db/internal/market_data"
	"github.com/achristie/save2db/pkg/cli"
	"github.com/achristie/save2db/pkg/platts"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	_ "modernc.org/sqlite"
)

var faCmd = &cobra.Command{
	Use:   "assessments",
	Short: "Assessment Data (Prices, Volumes, etc..)",
	Run: func(cmd *cobra.Command, args []string) {

		// create a platts api client
		client := platts.NewClient(viper.GetString("apikey"), viper.GetString("username"), viper.GetString("password"))

		// initialize DB and create necessary tables
		db := MD.NewDb("database2.db")
		as := MD.NewAssessmentsStore(db)

		// initial parameters
		start, err := time.Parse("2006-01-02T15:04:05", viper.GetString("startDate"))
		if err != nil {
			log.Fatal("Could not parse time: ", err)
		}

		p := cli.NewProgram(fmt.Sprintf("MDC: [%s], Modified Date >= [%s]", viper.GetString("mdc"), start), []string{"Assessments", "Deletes"})

		go func() {
			getAssessments(client, as, viper.GetString("mdc"), start, 10000, p)
			getDeletes(client, as, start, 10000, p)
		}()
		p.Start()
	},
}

func init() {
	fetchCmd.AddCommand(faCmd)
}

func getAssessments(client *platts.Client, db *MD.AssessmentsStore, MDC string, start time.Time, pageSize int, p *tea.Program) {
	data := make(chan platts.Result[platts.SymbolHistory])
	client.GetHistoryByMDC(MDC, start, pageSize, data)
	a := []platts.Assessment{}

	for result := range data {
		if result.Err != nil {
			log.Printf("Error! %s", result.Err)
		} else {
			res := result.Message
			pu := cli.ProgressUpdater{Name: "Assessments", Percent: 1 / float64(res.Metadata.TotalPages)}
			p.Send(pu)
			log.Printf("Assessment Data: %d records received from page [%d] in [%s] (%d total records).",
				len(res.Results), res.Metadata.Page, res.Metadata.QueryTime, res.Metadata.Count)

			a = append(a, res.Flatten()...)
		}
	}
	if err := db.Add(a); err != nil {
		log.Printf("Error inserting records: %s", err)
	}
	p.Send(cli.ProgressUpdater{Name: "Assessments", StatusMessage: fmt.Sprintf("Complete! Added [%d records] to [assessments]", len(a))})
}

// Get Deleted Assessments and remove from `assessments` table
func getDeletes(client *platts.Client, db *MD.AssessmentsStore, start time.Time, pageSize int, p *tea.Program) {
	data := make(chan platts.Result[platts.SymbolCorrection])
	client.GetDeletes(start, pageSize, data)
	a := []platts.Assessment{}

	for result := range data {
		if result.Err != nil {
			log.Printf("Error! %s", result.Err)
		} else {
			res := result.Message
			pu := cli.ProgressUpdater{Name: "Deletes", Percent: 1 / float64(res.Metadata.TotalPages)}
			p.Send(pu)
			log.Printf("Deletes: %d records received from page [%d] in [%s] (%d total records).",
				len(res.Results), res.Metadata.Page, res.Metadata.QueryTime, res.Metadata.Count)

			a = append(a, res.Flatten()...)
		}
	}
	if err := db.Remove(a); err != nil {
		log.Printf("Error removing records: %s", err)
	}
	p.Send(cli.ProgressUpdater{Name: "Deletes", StatusMessage: fmt.Sprintf("Complete! Removed [%d records] from [assessments]", len(a))})
}
