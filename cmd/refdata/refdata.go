package main

import (
	"flag"
	"log"

	RD "github.com/achristie/save2db/internal/ref_data"
	"github.com/achristie/save2db/pkg/platts"
	_ "github.com/mattn/go-sqlite3"
)

func main() {
	// read cmd line arguments
	APIKey := flag.String("apikey", "NULL", "API Key to call API with")
	Username := flag.String("username", "NULL", "Username to get a token")
	Password := flag.String("password", "NULL", "Password associated with Username")
	PageSize := flag.Int("p", 1000, "The page size to use for API Calls. Max is 10,000")
	flag.Parse()

	client := platts.NewClient(APIKey, Username, Password)

	db := RD.InitializeDb("database.db")

	page := 1
	for {
		rd, err := client.GetRefData(page, *PageSize)
		if err != nil {
			log.Println(err)
		}
		log.Printf("[%d] records received from page [%d] (%d total records). Adding to DB", len(rd.Results), rd.Metadata.Page, rd.Metadata.Count)
		if err := db.Add(rd); err != nil {
			log.Print(err)
		}
		if rd.Metadata.TotalPages == page || rd.Metadata.TotalPages == 0 {
			break
		}
		page++
	}

}
