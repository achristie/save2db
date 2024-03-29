package symbols_test

import (
	"path/filepath"
	"testing"

	"github.com/achristie/save2db/internal/services/symbols"
	"github.com/achristie/save2db/internal/sqlite"
	"github.com/achristie/save2db/pkg/platts"
)

func TestSymbolServiceAdd(t *testing.T) {
	t.Run("OK", func(t *testing.T) {
		db := sqlite.NewDB(filepath.Join(t.TempDir(), "db"))
		if err := db.Open(); err != nil {
			t.Fatal(err)
		}

		as, err := symbols.New(db.GetDB(), "SQLite")
		if err != nil {
			t.Fatal(err)
		}

		tx, err := db.BeginTx(nil)
		if err != nil {
			t.Fatal(err)
		}
		defer tx.Rollback()
		r := platts.SymbolResults{Symbol: "ABC", SettlementType: "Physical"}

		expected := 1
		res, err := as.Add(tx, r)
		if err != nil {
			t.Fatal(err)
		}

		result, err := res.RowsAffected()
		if err != nil {
			t.Fatal(err)
		}
		if result != int64(expected) {
			t.Errorf("expected: %d, result: %d", expected, result)
		}

		tx.Commit()

	})
}
