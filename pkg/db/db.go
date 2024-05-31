package db

import (
	"database/sql"
	"fmt"
	"net/url"
	"os"

	"web-golang-101/sqlc/queries"
)

type DBC struct {
	DB *sql.DB
}

func NewConnection() (*DBC, error) {
	connStr := os.Getenv("DB_STRING")
	if connStr == "" {
		return nil, fmt.Errorf("DB_STRING environment variable is not set")
	}

	u, err := url.Parse(connStr)
	if err != nil {
		return nil, err
	}

	db, err := sql.Open(u.Scheme, connStr)
	if err != nil {
		return nil, fmt.Errorf("500 | %w", err)
	}

	return &DBC{DB: db}, nil
}

func (d *DBC) NewQuery() *queries.Queries {
	return queries.New(d.DB)
}
