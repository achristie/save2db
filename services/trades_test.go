package services_test

import (
	"context"
	"path/filepath"
	"testing"

	"github.com/achristie/save2db/pkg/platts"
	"github.com/achristie/save2db/services"
	"github.com/achristie/save2db/sqlite"
)

func TestTradeSericeAdd(t *testing.T) {
	t.Run("SQLite", func(t *testing.T) {
		db := sqlite.NewDB(filepath.Join(t.TempDir(), "db"))
		if err := db.Open(); err != nil {
			t.Fatal(err)
		}

		ctx := context.Background()

		ts, err := services.NewTradeService(ctx, db.GetDB(), "SQLite")
		if err != nil {
			t.Fatal(err)
		}

		tx, err := db.BeginTx(ctx, nil)
		if err != nil {
			t.Fatal(err)
		}
		defer tx.Rollback()
		r := platts.TradeResults{Market: []string{"EU BFOE", "Other"}, DealID: 100}

		expected := 1
		res, err := ts.Add(ctx, tx, r)
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
