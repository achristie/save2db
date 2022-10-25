package trades_test

import (
	"path/filepath"
	"testing"

	"github.com/achristie/save2db/internal/services/trades"
	"github.com/achristie/save2db/internal/sqlite"
	"github.com/achristie/save2db/pkg/platts"
)

func TestTradeSericeAdd(t *testing.T) {
	t.Run("SQLite", func(t *testing.T) {
		db := sqlite.NewDB(filepath.Join(t.TempDir(), "db"))
		if err := db.Open(); err != nil {
			t.Fatal(err)
		}

		ts, err := trades.New(db.GetDB(), "SQLite")
		if err != nil {
			t.Fatal(err)
		}

		tx, err := db.BeginTx(nil)
		if err != nil {
			t.Fatal(err)
		}
		defer tx.Rollback()
		r := platts.TradeResults{Market: []string{"EU BFOE", "Other"}, DealID: 100}

		expected := 1
		res, err := ts.Add(tx, r)
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
