package assessments_test

import (
	"fmt"
	"path/filepath"
	"testing"

	"github.com/achristie/save2db/internal/services/assessments"
	"github.com/achristie/save2db/internal/sqlite"
	"github.com/achristie/save2db/pkg/platts"
)

func TestAssessmentsService_Add(t *testing.T) {
	t.Run("OK", func(t *testing.T) {
		db := sqlite.NewDB(filepath.Join(t.TempDir(), "db"))
		if err := db.Open(); err != nil {
			t.Fatal(err)
		}

		as, err := assessments.New(db.GetDB(), "SQLite")
		if err != nil {
			t.Fatal(err)
		}

		tx, err := db.BeginTx(nil)
		if err != nil {
			t.Fatal(err)
		}
		defer tx.Rollback()
		r := platts.Assessment{Symbol: "A", Bate: "B", AssessDate: "2022-01-01T00:00:00", ModDate: "2021-01-01T00:00:00", IsCorrected: "N", Value: 100}

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
	t.Run("Stress", func(t *testing.T) {
		db := sqlite.NewDB(filepath.Join(t.TempDir(), "db"))
		// db := sqlite.NewDB("test.db")
		if err := db.Open(); err != nil {
			t.Fatal(err)
		}

		as, err := assessments.New(db.GetDB(), "SQLite")
		if err != nil {
			t.Fatal(err)
		}

		tx, err := db.BeginTx(nil)
		if err != nil {
			t.Fatal(err)
		}
		defer tx.Rollback()
		for i := 0; i < 1000; i++ {
			r := platts.Assessment{Symbol: fmt.Sprint(i), Bate: "B", AssessDate: "2022-01-01T00:00:00", ModDate: "2021-01-01T00:00:00", IsCorrected: "N", Value: 100}
			_, err := as.Add(tx, r)
			if err != nil {
				t.Fatal(err)
			}
		}
		tx.Commit()
	})
}
func TestAssessmentsService_Remove(t *testing.T) {
	t.Run("SQLite", func(t *testing.T) {
		db := sqlite.NewDB(filepath.Join(t.TempDir(), "db"))
		if err := db.Open(); err != nil {
			t.Fatal(err)
		}

		as, err := assessments.New(db.GetDB(), "SQLite")
		if err != nil {
			t.Fatal(err)
		}

		tx, err := db.BeginTx(nil)
		if err != nil {
			t.Fatal(err)
		}
		// defer tx.Rollback()
		r := platts.Assessment{Symbol: "A", Bate: "B", AssessDate: "2022-01-01T00:00:00", ModDate: "2021-01-01T00:00:00", IsCorrected: "N", Value: 100}

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

		tx, err = db.BeginTx(nil)
		if err != nil {
			t.Fatal(err)
		}

		res, err = as.Remove(tx, r)
		fmt.Print(res.RowsAffected())

	})
}
