package main

import (
	"flag"
	"log"
	"os"
	"time"

	MD "github.com/achristie/save2db/internal/market_data"
	"github.com/achristie/save2db/pkg/cli"
	platts "github.com/achristie/save2db/pkg/platts"
	tea "github.com/charmbracelet/bubbletea"

	_ "modernc.org/sqlite"
)

func main() {
	// read cmd line arguments
	APIKey := flag.String("apikey", "NULL", "API Key to call API with")
	Username := flag.String("username", "NULL", "Username to get a token")
	Password := flag.String("password", "NULL", "Password associated with Username")
	StartDate := flag.String("t", time.Now().UTC().AddDate(0, 0, -3).Format("2006-01-02T15:04:05"), "Get updates since date. Format 2006-01-02T15:04:05")
	PageSize := flag.Int("p", 5000, "The page size to use for API Calls. Max is 10,000")
	MDC := flag.String("mdc", "NULL", "The MDC (Market Data Category) to fetch data for")
	Type := flag.String("type", "E", "Type of data to fetch. A - Assessments, D - Deleted Assessments, S - Symbol Data, E - Everything (A, D, S)")
	flag.Parse()

	f, err := os.OpenFile("log.log", os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		log.Fatalf("error opening file: %v", err)
	}
	defer f.Close()
	log.SetOutput(f)

	// create a platts api client
	client := platts.NewClient(APIKey, Username, Password)

	// initialize DB and create necessary tables
	db := MD.NewDb("database.db")
	as := MD.NewAssessmentsStore(db)
	rs := MD.NewSymbolStore(db)

	// initial parameters
	start, err := time.Parse("2006-01-02T15:04:05", *StartDate)
	if err != nil {
		log.Fatal("Could not parse time", err)
	}

	p := cli.NewProgram([]string{"Assessments", "Symbols", "Deletes"})

	go func() {
		// All or History
		if *Type == "E" || *Type == "A" {
			GetAssessments(client, as, *MDC, start, *PageSize, p)
		}

		// All or Reference
		if *Type == "E" || *Type == "S" {
			GetReferenceData(client, rs, start, *PageSize, p)
		}

		// All or Deletes
		if *Type == "E" || *Type == "D" {
			GetDeletes(client, as, start, *PageSize, p)
		}
	}()
	p.Start()
}

// Get Price Assessments and put into `assessments` table
func GetAssessments(client *platts.Client, db *MD.AssessmentsStore, MDC string, start time.Time, pageSize int, p *tea.Program) {
	data := make(chan platts.Result[platts.SymbolHistory])
	client.GetHistoryByMDC(MDC, start, pageSize, data)

	for result := range data {
		if result.Err != nil {
			log.Printf("Error! %s", result.Err)
		} else {
			res := result.Message
			pu := cli.ProgressUpdater{Name: "Assessments", Percent: 1 / float64(res.Metadata.TotalPages)}
			p.Send(pu)
			log.Printf("Assessment Data: %d records received from page [%d] in [%s] (%d total records). Adding to DB",
				len(res.Results), res.Metadata.Page, res.Metadata.QueryTime, res.Metadata.Count)
			if err := db.Add(res); err != nil {
				log.Printf("Error inserting records: %s", err)
			}
		}
	}
}

// Get Reference Data and put into `ref_data` table
func GetReferenceData(client *platts.Client, db *MD.SymbolStore, start time.Time, pageSize int, p *tea.Program) {
	data := make(chan platts.Result[platts.SymbolData])
	client.GetReferenceData(start, pageSize, data)

	for result := range data {
		if result.Err != nil {
			log.Printf("Error! %s", result.Err)
		} else {
			res := result.Message
			pu := cli.ProgressUpdater{Name: "Symbols", Percent: 1 / float64(res.Metadata.TotalPages)}
			p.Send(pu)
			log.Printf("Reference Data: %d records received from page [%d] in [%s] (%d total records). Adding to DB",
				len(res.Results), res.Metadata.Page, res.Metadata.QueryTime, res.Metadata.Count)
			if err := db.Add(res); err != nil {
				log.Printf("Error inserting records: %s", err)
			}
		}
	}
}

// Get Deleted Assessments and remove from `assessments` table
func GetDeletes(client *platts.Client, db *MD.AssessmentsStore, start time.Time, pageSize int, p *tea.Program) {
	data := make(chan platts.Result[platts.SymbolCorrection])
	client.GetDeletes(start, pageSize, data)

	for result := range data {
		if result.Err != nil {
			log.Printf("Error! %s", result.Err)
		} else {
			res := result.Message
			pu := cli.ProgressUpdater{Name: "Deletes", Percent: 1 / float64(res.Metadata.TotalPages)}
			p.Send(pu)
			log.Printf("Deletes: %d records received from page [%d] in [%s] (%d total records). Removing from DB",
				len(res.Results), res.Metadata.Page, res.Metadata.QueryTime, res.Metadata.Count)
			if err := db.Remove(res); err != nil {
				log.Printf("Error removing records: %s", err)
			}
		}
	}
}
