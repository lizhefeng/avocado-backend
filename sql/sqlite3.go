package sql

import (
	// this is required for sqlite3 usage
	_ "github.com/mattn/go-sqlite3"
)

// NewSQLite3 instantiates an sqlite3
func NewSQLite3(dbPath string) Store {
	return newStoreBase("sqlite3", dbPath)
}
