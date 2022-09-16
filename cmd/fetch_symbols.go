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

var symCmd = &cobra.Command{
	Use:   "symbols",
	Short: "Fetch Symbol Reference Data",
	Run: func(cmd *cobra.Command, args []string) {
		log.Printf("%s, %s", viper.GetString("startDate"), viper.GetString("mdc"))
		// create a platts api client
		client := platts.NewClient(viper.GetString("apikey"), viper.GetString("username"), viper.GetString("password"))

		// initialize DB and create necessary tables
		db := MD.NewDb("database2.db")
		ss := MD.NewSymbolStore(db)

		// initial parameters
		start, err := time.Parse("2006-01-02T15:04:05", viper.GetString("startDate"))
		if err != nil {
			log.Fatal("Could not parse time: ", err)
		}

		p := cli.NewProgram(fmt.Sprintf("MDC: [%s], Modified Date >= [%s]", viper.GetString("mdc"), start), []string{"Symbols"})

		go func() {
			getReferenceData(client, ss, start, viper.GetString("mdc"), 1000, p)
		}()
		p.Start()
	},
}

func init() {
	fetchCmd.AddCommand(symCmd)
	// a := &config{}
	// viper.Unmarshal(a)
	// log.Printf("%v", a)
}

// Get Reference Data and put into `symbols` table
func getReferenceData(client *platts.Client, db *MD.SymbolStore, start time.Time, mdc string, pageSize int, p *tea.Program) {
	data := make(chan platts.Result[platts.SymbolData])
	client.GetReferenceData(start, pageSize, mdc, data)
	sr := []platts.SymbolResults{}

	for result := range data {
		if result.Err != nil {
			log.Printf("Error! %s", result.Err)
		} else {
			res := result.Message
			pu := cli.ProgressUpdater{Name: "Symbols", Percent: 1 / float64(res.Metadata.TotalPages)}
			p.Send(pu)
			log.Printf("Reference Data: %d records received from page [%d] in [%s] (%d total records).",
				len(res.Results), res.Metadata.Page, res.Metadata.QueryTime, res.Metadata.Count)

			// add to temp slice
			sr = append(sr, res.Results...)

		}
	}

	if err := db.Add(sr); err != nil {
		log.Printf("Error inserting records: %s", err)
	}
	fmt.Printf("Added [%d records] to [symbols]", len(sr))
}
