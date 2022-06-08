package main

import (
	"flag"
	"log"
	"time"

	"github.com/achristie/save2db"
	platts "github.com/achristie/save2db/pkg/platts"

	_ "github.com/mattn/go-sqlite3"
)

func main() {
	// read cmd line arguments
	APIKey := flag.String("apikey", "NULL", "API Key to call API with")
	Username := flag.String("username", "NULL", "Username to get a token")
	Password := flag.String("password", "NULL", "Password associated with Username")
	flag.Parse()

	// create a platts api client
	client := platts.NewClient(APIKey, Username, Password)

	sc, err := client.GetDeletes(time.Now().AddDate(0, -2, 0), 1, 1000)

	log.Printf("%+v", sc)

	// initialize DB and create market_data table if it does not exist
	MarketDataStore := save2db.InitializeDb("database10.db")

	i, err := MarketDataStore.Remove(sc)
	if err != nil {
		log.Printf("error deleting records: %s", err)
	}
	log.Printf("deleted [%d] records from DB", i)

	log.Fatal()

	// initial parameters
	page := 1
	MDCs, err := client.GetSubscribedMDC()
	if err != nil {
		log.Fatalf("Could not get list of MDCs: %s", err)
	}
	start := MarketDataStore.GetLatestOrDefaultModifiedDate()
	pageSize := 20

	// loop through every MDC
	for _, v := range MDCs {
		// loop until everything is fetched
		for {
			// call history endpoint
			sh, err := client.GetHistoryByMDC(v, start, page, pageSize)

			// if there is an error then log it and go to the next MDC
			if err != nil {
				log.Print(err)
				break
			}

			// add data to database
			i, err := MarketDataStore.Add(sh)
			if err != nil {
				log.Print("error inserting records: ", err)
			}
			log.Printf("added [%d] records to DB", i)

			// exit loop when all data has been fetched
			if sh.Metadata.TotalPages == page || sh.Metadata.TotalPages == 0 {
				break
			}

			// avoid getting throttled by the API
			time.Sleep(2 * time.Second)
			page += 1
		}
		// reset page counter, sleep for throttling
		page = 1
		time.Sleep(2 * time.Second)
	}
}
