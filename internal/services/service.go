package services

import (
	"context"
	"database/sql"
)

type Service interface {
	Add(context.Context, *sql.Tx, interface{}) (sql.Result, error)
	Remove(context.Context, *sql.Tx, interface{}) (sql.Result, error)
}
