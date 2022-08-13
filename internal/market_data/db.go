package market_data

import (
	"database/sql"
	"log"
	"os"
)

func NewDb(FileName string) *sql.DB {
	file, err := os.OpenFile(FileName, os.O_RDWR|os.O_CREATE, 0666)

	if err != nil {
		log.Fatalf("could not open %s %v", FileName, err)
	}

	db, err := sql.Open("sqlite3", file.Name())
	if err != nil {
		log.Fatal(err)
	}

	return db
}
