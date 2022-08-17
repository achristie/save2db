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
	Type := flag.String("type", "Ref", "placeholder")
	flag.Parse()

	// create a platts api client
	client := platts.NewClient(APIKey, Username, Password)

	// initialize DB and create necessary tables
	db := MD.NewDb("database.db")
	as := MD.NewAssessmentsStore(db)
	rs := MD.NewRefDataStore(db)

	// initial parameters
	start, err := time.Parse("2006-01-02T15:04:05", *StartDate)
	if err != nil {
		log.Fatal("Could not parse time", err)
	}

	if *Type == "Ref" {
		// Update `ref_data` table with any reference data modified since `t`
		GetReferenceData(client, rs, start, min(1000, *PageSize))

	} else {

		// Update `assessments` table with records modified since `t`
		GetAssessments(client, as, *MDC, start, min(10000, *PageSize))

		// Update `assessments` table with any corrections (deletes) since `t`
		GetCorrections(client, as, start, *PageSize)
	}

}

// Uses the `client` to fetch historical data for given MDC modified since `start`
// Uses the concurrent get history method to fetch data in parallel
// Store results in DB
func GetAssessments(client *platts.Client, db *MD.AssessmentsStore, MDC string, start time.Time, pageSize int) {
	ch := make(chan platts.Result[platts.SymbolHistory])

	go func() {
		log.Printf("Fetching history for [%s] since %s", MDC, start.String())
		client.GetHistoryByMDCConcurrent(MDC, start, pageSize, ch)
	}()

	for result := range ch {
		if result.Err != nil {
			log.Printf("Error retrieving data: %s", result.Err)
		} else {
			log.Printf("[%d] records received from page [%d] in [%s] (%d total records). Adding to DB",
				len(result.OK.Results), result.OK.Metadata.Page, result.OK.Metadata.QueryTime, result.OK.Metadata.Count)
			if err := db.Add(result.OK); err != nil {
				log.Printf("Error inserting records: %s", err)
			}
		}
	}
}

func GetReferenceData(client *platts.Client, db *MD.RefDataStore, start time.Time, pageSize int) {
	ch := make(chan interface{})
	client.GetRefData(start, pageSize, ch)

	for {
		select {
		// case err := <-errs:
		// 	// log.Println(err)

		case result := <-ch:
			res := result.(platts.ReferenceData)
			// log.Println(res.Metadata.TotalPages)
			log.Printf("%d records received from page [%d] in [%s] (%d total records). Removing from DB",
				len(res.Results), res.Metadata.Page, res.Metadata.QueryTime, res.Metadata.Count)

		}
	}

}

// ch := make(chan platts.Result[platts.ReferenceData])

// go func() {
// 	log.Printf("Fetching Reference Data updated since %s", start.String())
// 	client.GetRefDataConcurrent(start, pageSize, ch)
// }()

// for result := range ch {
// 	if result.Err != nil {
// 		log.Printf("Error retrieving data: %s", result.Err)
// 	} else {
// 		log.Printf("%d records received from page [%d] in [%s] (%d total records). Removing from DB",
// 			len(result.OK.Results), result.OK.Metadata.Page, result.OK.Metadata.QueryTime, result.OK.Metadata.Count)
// 		if err := db.Add(result.OK); err != nil {
// 			log.Printf("Error updating records: %s", err)
// 		}
// 	}
// }

func GetCorrections(client *platts.Client, db *MD.AssessmentsStore, start time.Time, pageSize int) {
	ch := make(chan platts.Result[platts.SymbolCorrection])

	go func() {
		log.Printf("Fetching corrections since %s", start.String())
		err := client.GetDeletesConcurrent(start, pageSize, ch)
		if err != nil {
			log.Fatal(err)
		}
	}()

	for result := range ch {
		if result.Err != nil {
			log.Printf("Error retrieving data: %s", result.Err)
		} else {
			log.Printf("%d records received from page [%d] in [%s] (%d total records). Removing from DB",
				len(result.OK.Results), result.OK.Metadata.Page, result.OK.Metadata.QueryTime, result.OK.Metadata.Count)
			if err := db.Remove(result.OK); err != nil {
				log.Printf("Error removing records: %s", err)
			}
		}
	}
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}
