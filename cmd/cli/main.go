package main

import (
	"database/sql"
	"flag"
	"log"
	"os"

	"github.com/achristie/save2db/pkg/plattsapi"

	_ "github.com/mattn/go-sqlite3"
)

func main() {
	APIKey := flag.String("apikey", "NULL", "API Key to call API with")
	Username := flag.String("username", "NULL", "Username to get a token")
	Password := flag.String("password", "NULL", "Password associated with Username")

	file, err := os.Create("database.db")
	if err != nil {
		log.Fatal(err)
	}

	file.Close()

	db, err := sql.Open("sqlite3", file.Name())
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	createTable(db)
	if err != nil {
		log.Fatal(err)
	}
	// addRecord(db, "ABC", "B", "N", time.Now(), time.Now(), 100.12)

	flag.Parse()
	c := plattsapi.New(APIKey, Username, Password)
	data := c.CallApi()

	// log.Printf("%+v", data)
	for _, v := range data.Results {
		for _, v2 := range v.Data {
			addRecord(db, v.Symbol, v2.Bate, v2.IsCorrected, v2.ModDate, v2.AssessDate, v2.Value)
		}
	}

	// t := GetToken(*Username, *Password, *APIKey)

	// log.Print(t.AccessToken)
	// log.Println(t.RefreshToken)
}

func addRecord(db *sql.DB, Symbol string, Bate string, IsCorrected string, ModifiedDate string, AssessedDate string, Price float32) {
	records := `INSERT INTO market_data(symbol, bate, price, assessed_date, modified_date, is_corrected) VALUES(?, ?, ?, ?, ?, ?)`
	query, err := db.Prepare(records)
	if err != nil {
		log.Fatal(err)
	}
	_, err = query.Exec(Symbol, Bate, Price, AssessedDate, ModifiedDate, IsCorrected)
	if err != nil {
		log.Fatal(err)
	}
}

func createTable(db *sql.DB) {
	market_data_table := `CREATE TABLE market_data (
		"symbol" TEXT NOT NULL,
		"bate" TEXT NOT NULL,
		"price" NUM NOT NULL,
		"assessed_date" datetime NOT NULL,
		"modified_date" datetime NOT NULL,
		"is_corrected" string NOT NULL,
		PRIMARY KEY (symbol, bate, assessed_date) );`
	query, err := db.Prepare(market_data_table)
	if err != nil {
		log.Fatal(err)
	}
	query.Exec()
	log.Println("Market Data table created succesfully")
}
