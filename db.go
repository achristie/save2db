package save2db

import (
	"database/sql"
	"log"
	"os"
	"time"

	platts "github.com/achristie/save2db/pkg/platts"
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
		"is_corrected" string NOT NULL
	);`
	// PRIMARY KEY (symbol, bate, assessed_date) );`
	query, err := db.Prepare(market_data_table)
	if err != nil {
		log.Fatal(err)
	}
	query.Exec()
	log.Println("db: market_data table created succesfully")
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

func (m *MarketDataStore) GetLatestOrDefaultModifiedDate() time.Time {
	// must cast into a string because of gosql driver issues
	row := m.database.QueryRow("SELECT CAST(max(modified_date) as text) from market_data")

	var result sql.NullString
	err := row.Scan(&result)
	defDate := time.Now().UTC().AddDate(0, 0, -7)

	if err != nil {
		return defDate
	}
	t, err := time.Parse("2006-01-02T15:04:05", result.String)
	if err != nil {
		log.Printf("db: Modified Date is null. Returning Default (Now - 7 Days) Value: %s", defDate)
		return defDate
	}
	return t
}

func (m *MarketDataStore) Remove(data platts.SymbolCorrection) (int, error) {
	var records []dbClass
	for _, v := range data.Results {
		for _, v2 := range v.Data {
			records = append(records, dbClass{
				symbol:       v.Symbol,
				bate:         v2.Bate,
				assessedDate: v2.AssessDate,
			})
		}
	}
	err := m.del(records)
	if err != nil {
		return 0, err
	}
	return len(records), nil
}

func (m *MarketDataStore) del(records []dbClass) error {
	del := `DELETE FROM market_data where symbol = ? and bate = ? and assessed_date = ?`
	query, err := m.database.Prepare(del)
	if err != nil {
		return err
	}
	defer query.Close()

	// bulk delete
	tx, err := m.database.Begin()
	if err != nil {
		return err
	}

	for _, r := range records {
		_, err := tx.Stmt(query).Exec(r.symbol, r.bate, r.assessedDate)
		if err != nil {
			tx.Rollback()
			return err
		}
	}
	return tx.Commit()
}

func (m *MarketDataStore) Add(data platts.SymbolHistory) (int, error) {
	var records []dbClass
	// change data structure for ease of INSERT
	for _, v := range data.Results {
		for _, v2 := range v.Data {
			records = append(records, dbClass{
				symbol:       v.Symbol,
				bate:         v2.Bate,
				price:        v2.Value,
				modifiedDate: v2.ModDate,
				assessedDate: v2.AssessDate,
				isCorrected:  v2.IsCorrected,
			})
		}
	}
	err := m.insert(records)
	if err != nil {
		return 0, err
	}
	return len(records), nil
}

func (m *MarketDataStore) insert(records []dbClass) error {
	ins := `INSERT or REPLACE INTO market_data(symbol, bate, price, assessed_date, modified_date, is_corrected) VALUES(?, ?, ?, ?, ?, ?)`
	query, err := m.database.Prepare(ins)
	if err != nil {
		return err
	}
	defer query.Close()

	// bulk insert
	tx, err := m.database.Begin()
	if err != nil {
		return err
	}

	for _, r := range records {
		_, err := tx.Stmt(query).Exec(r.symbol, r.bate, r.price, r.assessedDate, r.modifiedDate, r.isCorrected)
		if err != nil {
			tx.Rollback()
			return err
		}
	}
	return tx.Commit()
}
