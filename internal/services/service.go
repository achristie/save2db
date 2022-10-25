package services

import (
	"database/sql"
)

type Service interface {
	Add(*sql.Tx, interface{}) (sql.Result, error)
	Remove(*sql.Tx, interface{}) (sql.Result, error)
}
