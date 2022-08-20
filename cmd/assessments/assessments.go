package main

import (
	"flag"
	"log"
	"time"

	MD "github.com/achristie/save2db/internal/market_data"
	platts "github.com/achristie/save2db/pkg/platts"

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
	Type := flag.String("type", "A", "Type of data to fetch. H - HistoricalAssessments, D - AssessmentsDeleted, R - ReferenceData, A - All (H, D, R)")
	flag.Parse()

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

	// All or History
	if *Type == "A" || *Type == "H" {
		GetAssessments(client, as, *MDC, start, *PageSize)
	}

	// All or Reference
	if *Type == "A" || *Type == "R" {
		GetReferenceData(client, rs, start, *PageSize)
	}

	// All or Deletes
	if *Type == "A" || *Type == "D" {
		GetDeletes(client, as, start, *PageSize)
	}
}

// Get Price Assessments and put into `assessments` table
func GetAssessments(client *platts.Client, db *MD.AssessmentsStore, MDC string, start time.Time, pageSize int) {
	data := make(chan platts.Result[platts.SymbolHistory])
	client.GetHistoryByMDC(MDC, start, pageSize, data)

	for result := range data {
		if result.Err != nil {
			log.Printf("Error! %s", result.Err)
		} else {
			res := result.Message
			log.Printf("Assessment Data: %d records received from page [%d] in [%s] (%d total records). Adding to DB",
				len(res.Results), res.Metadata.Page, res.Metadata.QueryTime, res.Metadata.Count)
			if err := db.Add(res); err != nil {
				log.Printf("Error inserting records: %s", err)
			}
		}
	}
}

// Get Reference Data and put into `ref_data` table
func GetReferenceData(client *platts.Client, db *MD.SymbolStore, start time.Time, pageSize int) {
	data := make(chan platts.Result[platts.SymbolData])
	client.GetReferenceData(start, pageSize, data)

	for result := range data {
		if result.Err != nil {
			log.Printf("Error! %s", result.Err)
		} else {
			res := result.Message
			log.Printf("Reference Data: %d records received from page [%d] in [%s] (%d total records). Adding to DB",
				len(res.Results), res.Metadata.Page, res.Metadata.QueryTime, res.Metadata.Count)
			if err := db.Add(res); err != nil {
				log.Printf("Error inserting records: %s", err)
			}
		}
	}
}

// Get Deleted Assessments and remove from `assessments` table
func GetDeletes(client *platts.Client, db *MD.AssessmentsStore, start time.Time, pageSize int) {
	data := make(chan platts.Result[platts.SymbolCorrection])
	client.GetDeletes(start, pageSize, data)

	for result := range data {
		if result.Err != nil {
			log.Printf("Error! %s", result.Err)
		} else {
			res := result.Message
			log.Printf("Deletes: %d records received from page [%d] in [%s] (%d total records). Removing from DB",
				len(res.Results), res.Metadata.Page, res.Metadata.QueryTime, res.Metadata.Count)
			if err := db.Remove(res); err != nil {
				log.Printf("Error removing records: %s", err)
			}
		}
	}
}
