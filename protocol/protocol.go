package protocol

import (
	"context"
	"github.com/pkg/errors"
)

const (
	// ReviewTableName is the table name of review history
	ReviewTableName = "review"
)

var (
	// ErrNotExist indicates certain item does not exist in Blockchain database
	ErrNotExist = errors.New("not exist in DB")
	// ErrAlreadyExist indicates certain item already exists in Blockchain database
	ErrAlreadyExist = errors.New("already exist in DB")
	// ErrUnimplemented indicates a method is not implemented yet
	ErrUnimplemented = errors.New("method is unimplemented")
)

// Protocol defines the protocol interfaces for key store
type Protocol interface {
	CreateTables(context.Context) error
	ReviewHandler
}

// ReviewHandler ishte interface of handling reviews
type ReviewHandler interface {
	PutReview(context.Context, *ReviewInput) (*Review, error)
	GetReview(uint64) (*Review, error)
}

// Review defines the structure of reviews
type Review struct {
	ReviewID  uint64
	UserID    uint64
	ArtID     uint64
	Text      string
	TimeStamp uint64
	Upvotes   uint64
}

// ReviewInput defines the structure of review input
type ReviewInput struct {
	UserID    uint64
	ArtID     uint64
	Text      string
	TimeStamp uint64
	Upvotes   uint64
}
