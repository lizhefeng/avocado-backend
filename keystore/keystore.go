package keystore

import (
	"context"

	"github.com/pkg/errors"

	"github.com/lizhefeng/avocado-backend/protocol"
	"github.com/lizhefeng/avocado-backend/protocol/albums"
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
	if err := ks.store.Start(ctx); err != nil {
		return errors.Wrap(err, "failed to start db")
	}
	return ks.CreateTablesIfNotExist()
}

// Stop stops the indexer
func (ks *KeyStore) Stop(ctx context.Context) error {
	return ks.store.Stop(ctx)
}

// CreateTablesIfNotExist creates tables in local database
func (ks *KeyStore) CreateTablesIfNotExist() error {
	for _, p := range ks.registry.All() {
		if err := p.CreateTables(context.Background()); err != nil {
			return errors.Wrap(err, "failed to create a table")
		}
	}
	return nil
}

// RegisterProtocol registers a protocol to the indexer
func (ks *KeyStore) RegisterProtocol(protocolID string, protocol protocol.Protocol) error {
	return ks.registry.Register(protocolID, protocol)
}

// RegisterDefaultProtocols registers default protocols to the keystore
func (ks *KeyStore) RegisterDefaultProtocols() error {
	albumsProtocol := albums.NewProtocol(ks.store)
	return ks.RegisterProtocol(albums.ProtocolID, albumsProtocol)
}

// PutReview inserts a new review into the underlying db
func (ks *KeyStore) PutReview(ctx context.Context, protocolID string, reviewInput *protocol.ReviewInput) (*protocol.Review, error) {
	p, ok := ks.registry.Find(protocolID)
	if !ok {
		return nil, errors.New("protocol is unregistered")
	}
	return p.PutReview(ctx, reviewInput)
}

// GetReview gets a review from the underlying db
func (ks *KeyStore) GetReview(protocolID string, reviewID uint64) (*protocol.Review, error) {
	p, ok := ks.registry.Find(protocolID)
	if !ok {
		return nil, errors.New("protocol is unregistered")
	}
	return p.GetReview(reviewID)
}
