package cmd

import (
	"fmt"
	"log"
	"os"
	"time"

	MD "github.com/achristie/save2db/internal/market_data"
	"github.com/achristie/save2db/pkg/platts"
	tui "github.com/achristie/save2db/pkg/tui/progress"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/spf13/cobra"
	_ "modernc.org/sqlite"
)

var symCmd = &cobra.Command{
	Use:   "symbols",
	Short: "Fetch symbol reference data",
	Run: func(cmd *cobra.Command, args []string) {
		// // create a platts api client
		// client := platts.NewClient(viper.GetString("apikey"), viper.GetString("username"), viper.GetString("password"))

		// // initialize DB and create necessary tables
		// db := MD.NewDb("database2.db")
		// ss := MD.NewSymbolStore(db)

		// p := cli.NewProgram(fmt.Sprintf("MDC: [%s], Modified Date >= [%s]", mdc, startDate), []string{"Symbols"})

		// go func() {
		// 	getReferenceData(client, ss, startDate, mdc, 1000, p)
		// }()
		// p.Start()
	},
}

func init() {
	fetchCmd.AddCommand(symCmd)
}

// Get Reference Data and put into `symbols` table
func getReferenceData(client *platts.Client, db *MD.SymbolStore, start time.Time, mdc string, pageSize int, p *tea.Program) {
	data := make(chan platts.Result[platts.SymbolData])
	client.GetReferenceData(start, pageSize, mdc, data)
	sr := []platts.SymbolResults{}

	for result := range data {
		if result.Err != nil {
			log.Printf("Error - Reference Data:  %s", result.Err)
			p.Send(tui.StatusUpdater{Name: "Symbols", Status: tui.Status{Category: tui.ERROR, Msg: "An error occured, please retry."}})
			os.Exit(1)
		} else {
			res := result.Message
			pu := tui.ProgressUpdater{Name: "Symbols", Percent: 1 / float64(res.Metadata.TotalPages)}
			p.Send(pu)
			log.Printf("Reference Data: %d records received from page [%d] in [%s] (%d total records).",
				len(res.Results), res.Metadata.Page, res.Metadata.QueryTime, res.Metadata.Count)

			// add to temp slice
			sr = append(sr, res.Results...)

		}
	}

	if err := db.Add(sr); err != nil {
		log.Printf("Error inserting records: %s", err)
		p.Send(tui.StatusUpdater{Name: "Symbols", Status: tui.Status{Category: tui.ERROR, Msg: "An error occured, please retry."}})
	}
	p.Send(tui.StatusUpdater{Name: "Symbols", Status: tui.Status{Category: tui.COMPLETED, Msg: fmt.Sprintf("Complete! Added [%d records] to [symbols]", len(sr))}})
}
