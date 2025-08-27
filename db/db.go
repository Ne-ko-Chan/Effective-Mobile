package db

import (
	"database/sql"

	_ "github.com/lib/pq"
)

func NewPSQLStorage(connStr string) (*sql.DB, error) {
	db,err := sql.Open("postgres",connStr)
	if err != nil {
		return nil, err
	}

	return db, nil
}
