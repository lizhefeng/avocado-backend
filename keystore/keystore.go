package keystore

import (
	"context"
	"database/sql"

	"github.com/pkg/errors"

	"github.com/lizhefeng/avocado-backend/protocol"
	s "github.com/lizhefeng/avocado-backend/sql"
)

// KeyStore handles front end requests
type KeyStore struct {
	store    s.Store
	registry *protocol.Registry
}

// Config contains indexer configs
type Config struct {
	DBPath string `yaml:"dbPath"`
}

// NewKeyStore creates a new keystore
func NewKeyStore(store s.Store) *KeyStore {
	return &KeyStore{
		store:    store,
		registry: &protocol.Registry{},
	}
}

// Start starts the keystore
func (ks *KeyStore) Start(ctx context.Context) error {
	return ks.store.Start(ctx)
}

// Stop stops the indexer
func (ks *KeyStore) Stop(ctx context.Context) error {
	return ks.store.Stop(ctx)
}

// PutReview inserts a new review into the underlying db
func (ks *KeyStore) PutReview(ctx context.Context, protocolID string, review *protocol.Review) error {
	p, ok := ks.registry.Find(protocolID)
	if !ok {
		return errors.New("protocol is unregistered")
	}
	return ks.store.Transact(func(tx *sql.Tx) error {
		return p.PutReview(ctx, tx, review)
	})
}
