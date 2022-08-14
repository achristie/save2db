package main

import (
	"flag"
	"log"
	"time"

	MD "github.com/achristie/save2db/internal/market_data"
	platts "github.com/achristie/save2db/pkg/platts"

	_ "github.com/mattn/go-sqlite3"
)

func main() {
	// read cmd line arguments
	APIKey := flag.String("apikey", "NULL", "API Key to call API with")
	Username := flag.String("username", "NULL", "Username to get a token")
	Password := flag.String("password", "NULL", "Password associated with Username")
	StartDate := flag.String("t", time.Now().UTC().AddDate(0, 0, -3).Format("2006-01-02T15:04:05"), "Get updates since date. Format 2006-01-02T15:04:05")
	PageSize := flag.Int("p", 5000, "The page size to use for API Calls. Max is 10,000")
	MDC := flag.String("mdc", "NULL", "The MDC (Market Data Category) to fetch data for")
	flag.Parse()

	// create a platts api client
	client := platts.NewClient(APIKey, Username, Password)

	// initialize DB and create market_data table if it does not exist
	db := MD.NewDb("database.db")
	as := MD.NewAssessmentsStore(db)

	// initial parameters
	start, err := time.Parse("2006-01-02T15:04:05", *StartDate)
	if err != nil {
		log.Fatal("Could not parse time", err)
	}

	// Update market_data table with records modified since `start`
	GetAssessments(client, as, *MDC, start, *PageSize)

}

// Uses the `client` to fetch historical data for given MDC modified since `start`
// Uses the concurrent get history method to fetch data in parallel
// Store results in DB
func GetAssessments(client *platts.Client, db *MD.AssessmentsStore, MDC string, start time.Time, pageSize int) {
	ch := make(chan platts.Result)

	go func() {
		log.Printf("Fetching history for [%s] since %s", MDC, start.String())
		err := client.GetHistoryByMDCConcurrent(MDC, start, pageSize, ch)
		if err != nil {
			log.Fatal(err)
		}
	}()

	for result := range ch {
		if result.Err != nil {
			log.Printf("Error retrieving data: %s", result.Err)
		} else {
			log.Printf("[%d] records received from page [%d] in [%s] (%d total records). Adding to DB",
				len(result.SH.Results), result.SH.Metadata.Page, result.SH.Metadata.QueryTime, result.SH.Metadata.Count)
			if err := db.Add(result.SH); err != nil {
				log.Printf("Error inserting records: %s", err)
			}
		}
	}
}
