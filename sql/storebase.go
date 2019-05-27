package sql

import (
	"context"
	"database/sql"
	"os"
	"sync"

	"github.com/rs/zerolog"
)

// Store is the interface of KV store.
type Store interface {
	Start(context.Context) error
	Stop(context.Context) error
	// Get DB instance
	GetDB() *sql.DB

	// Transact wrap the transaction
	Transact(func(*sql.Tx) error) error
}

// storebase is local sqlite3
type storeBase struct {
	mutex      sync.RWMutex
	db         *sql.DB
	connectStr string
	driverName string
}

// logger is initialized with default settings
var logger = zerolog.New(os.Stderr).Level(zerolog.InfoLevel).With().Timestamp().Logger()

// NewStoreBase instantiates an store base
func newStoreBase(driverName string, connectStr string) Store {
	return &storeBase{db: nil, connectStr: connectStr, driverName: driverName}
}

// Start opens the SQL (creates new file if not existing yet)
func (s *storeBase) Start(_ context.Context) error {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	if s.db != nil {
		return nil
	}

	// Use db to perform SQL operations on database
	db, err := sql.Open(s.driverName, s.connectStr)
	if err != nil {
		return err
	}
	s.db = db
	return nil
}

// Stop closes the SQL
func (s *storeBase) Stop(_ context.Context) error {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	if s.db != nil {
		err := s.db.Close()
		s.db = nil
		return err
	}
	return nil
}

// Stop closes the SQL
func (s *storeBase) GetDB() *sql.DB {
	s.mutex.RLock()
	defer s.mutex.RUnlock()

	return s.db
}

// Transact wrap the transaction
func (s *storeBase) Transact(txFunc func(*sql.Tx) error) (err error) {
	tx, err := s.db.Begin()
	if err != nil {
		return err
	}
	defer func() {
		switch {
		case recover() != nil:
			if rollbackErr := tx.Rollback(); rollbackErr != nil {
				logger.Error().Err(rollbackErr) // log err after Rollback
			}
		case err != nil:
			// err is non-nil; don't change it
			if rollbackErr := tx.Rollback(); rollbackErr != nil {
				logger.Error().Err(rollbackErr)
			}
		default:
			// err is nil; if Commit returns error update err
			if commitErr := tx.Commit(); commitErr != nil {
				logger.Error().Err(commitErr)
			}
		}
	}()
	err = txFunc(tx)
	return err
}
