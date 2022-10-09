package pg_test

import (
	"testing"

	"github.com/achristie/save2db/pg"
	_ "github.com/jackc/pgx/v5/stdlib"
)

func TestOpen(t *testing.T) {
	t.Run("Ok", func(t *testing.T) {
		db := pg.NewDB("postgres://postgres:password@localhost:5432/testdb")
		if err := db.Open(); err != nil {
			t.Fatal(err)
		}

	})
}
