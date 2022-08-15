package main

import (
	"flag"
	"log"
	"time"

	MD "github.com/achristie/save2db/internal/market_data"
	"github.com/achristie/save2db/pkg/platts"

	// _ "github.com/mattn/go-sqlite3"
	_ "modernc.org/sqlite"
)

func main() {
	// read cmd line arguments
	APIKey := flag.String("apikey", "NULL", "API Key to call API with")
	Username := flag.String("username", "NULL", "Username to get a token")
	Password := flag.String("password", "NULL", "Password associated with Username")
	PageSize := flag.Int("p", 1000, "The page size to use for API Calls. Max is 1,000")
	flag.Parse()

	client := platts.NewClient(APIKey, Username, Password)

	db := MD.NewDb("database.db")
	rds := MD.NewRefDataStore(db)

	page := 1
	for {
		rd, err := client.GetRefData(page, *PageSize)
		if err != nil {
			log.Println(err)
			break
		}
		log.Printf("[%d] records received from page [%d]  (%d total records) in [%s]ms. Adding to DB", len(rd.Results), rd.Metadata.Page, rd.Metadata.Count, rd.Metadata.QueryTime)
		if err := rds.Add(rd); err != nil {
			log.Print(err)
			break
		}
		if rd.Metadata.TotalPages == page || rd.Metadata.TotalPages == 0 {
			break
		}
		page++
		time.Sleep(500 * time.Millisecond)
	}

}
