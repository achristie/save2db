package save2db

import (
	"database/sql"
	"log"
	"os"

	"github.com/achristie/save2db/pkg/plattsapi"
)

type MarketDataStore struct {
	database *sql.DB
}

type dbClass struct {
	symbol       string
	bate         string
	modifiedDate string
	assessedDate string
	price        float32
	isCorrected  string
}

func createTable(db *sql.DB) {
	market_data_table := `CREATE TABLE IF NOT EXISTS market_data (
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
	log.Println("market_data table created succesfully")
}

func InitializeDb(dbFileName string) *MarketDataStore {
	file, err := os.OpenFile(dbFileName, os.O_RDWR|os.O_CREATE, 0666)

	if err != nil {
		log.Fatalf("could not open %s %v", dbFileName, err)
	}

	db, err := sql.Open("sqlite3", file.Name())
	if err != nil {
		log.Fatal(err)
	}

	createTable(db)

	return &MarketDataStore{database: db}
}

func (f *MarketDataStore) AddPricingData(data plattsapi.SymbolHistory) {
	for _, v := range data.Results {
		for _, v2 := range v.Data {
			record := dbClass{
				symbol:       v.Symbol,
				bate:         v2.Bate,
				price:        v2.Value,
				modifiedDate: v2.ModDate,
				assessedDate: v2.AssessDate,
				isCorrected:  v2.IsCorrected,
			}
			f.insert(&record)
		}
	}
}

func (f *MarketDataStore) insert(record *dbClass) {
	ins := `INSERT or REPLACE INTO market_data(symbol, bate, price, assessed_date, modified_date, is_corrected) VALUES(?, ?, ?, ?, ?, ?)`
	query, err := f.database.Prepare(ins)
	if err != nil {
		log.Fatal(err)
	}
	_, err = query.Exec(record.symbol, record.bate, record.price, record.assessedDate, record.modifiedDate, record.isCorrected)
	if err != nil {
		log.Println(err)
	}

}
