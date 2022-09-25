package pg_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/achristie/save2db/pg"
	"github.com/achristie/save2db/pkg/platts"
)

func TestAssessmentsService_Add(t *testing.T) {
	t.Run("OK", func(t *testing.T) {
		db := pg.NewDB("postgres://postgres:password@localhost:5432/testdb")
		if err := db.Open(); err != nil {
			t.Fatal(err)
		}

		ctx := context.Background()

		as, err := pg.NewAssessmentsService(ctx, db)
		if err != nil {
			t.Fatal(err)
		}

		tx, err := db.BeginTx(ctx, nil)
		if err != nil {
			t.Fatal(err)
		}
		defer tx.Rollback()
		r := platts.Assessment{Symbol: "A", Bate: "B", AssessDate: "2022-01-01T00:00:00", ModDate: "2021-01-01T00:00:00", IsCorrected: "N", Value: 100}

		expected := 1
		res, err := as.Add(ctx, tx, r)
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
		db := pg.NewDB("postgres://postgres:password@localhost:5432/testdb")
		if err := db.Open(); err != nil {
			t.Fatal(err)
		}

		ctx := context.Background()
		as, err := pg.NewAssessmentsService(ctx, db)
		if err != nil {
			t.Fatal(err)
		}

		tx, err := db.BeginTx(ctx, nil)
		if err != nil {
			t.Fatal(err)
		}
		defer tx.Rollback()
		for i := 0; i < 10000; i++ {
			r := platts.Assessment{Symbol: fmt.Sprint(i), Bate: "B", AssessDate: "2022-01-01T00:00:00", ModDate: "2021-01-01T00:00:00", IsCorrected: "N", Value: 100}
			_, err := as.Add(ctx, tx, r)
			if err != nil {
				t.Fatal(err)
			}
		}
		tx.Commit()
	})
}
