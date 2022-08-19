package main

import (
	"flag"
	"log"
	"time"

	TD "github.com/achristie/save2db/internal/trade_data"
	"github.com/achristie/save2db/pkg/platts"
	_ "modernc.org/sqlite"
)

func main() {
	// read cmd line arguments
	APIKey := flag.String("apikey", "NULL", "API Key to call API with")
	Username := flag.String("username", "NULL", "Username to get a token")
	Password := flag.String("password", "NULL", "Password associated with Username")
	StartDate := flag.String("t", time.Now().UTC().AddDate(0, 0, -3).Format("2006-01-02T15:04:05"), "Get updates since date. Format 2006-01-02T15:04:05")
	PageSize := flag.Int("p", 1000, "The page size to use for API Calls. Max is 10,000")
	flag.Parse()

	client := platts.NewClient(APIKey, Username, Password)

	db := TD.NewDb("database.db")
	tds := TD.NewTradeDataStore(db)

	// initial parameters
	start, err := time.Parse("2006-01-02T15:04:05", *StartDate)
	if err != nil {
		log.Fatal("Could not parse time", err)
	}

	GetTrades(client, tds, start, *PageSize)

}

func GetTrades(client *platts.Client, db *TD.TradeDataStore, start time.Time, pageSize int) {
	data := make(chan platts.Result[platts.TradeData])
	client.GetTradeData(start, pageSize, data)

	for result := range data {
		if result.Err != nil {
			log.Printf("Error %s", result.Err)
		} else {
			res := result.Message
			log.Printf("Trade Data: %d records received from page [%d] in [%s] (%d total records). Adding to DB",
				len(res.Results), res.Metadata.Page, res.Metadata.QueryTime, res.Metadata.Count)
			if err := db.Add(res); err != nil {
				log.Printf("Error inserting records: %s", err)
			}
		}
	}
}
