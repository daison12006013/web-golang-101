package db

import (
	"database/sql"
	"fmt"
	"net/url"
	"os"
	"strconv"

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
		return nil, err
	}

	maxConn, err := defaultMaxOpenConns()
	if err != nil {
		return nil, err
	}
	db.SetMaxOpenConns(*maxConn)

	return &DBC{DB: db}, nil
}

func (d *DBC) NewQuery() *queries.Queries {
	return queries.New(d.DB)
}

func defaultMaxOpenConns() (*int, error) {
	maxConn := os.Getenv("DB_MAX_OPEN_CONNS")
	if maxConn == "" {
		maxConn = "20"
	}

	maxConnInt, err := strconv.Atoi(maxConn)
	if err != nil {
		return nil, err
	}

	return &maxConnInt, nil
}
