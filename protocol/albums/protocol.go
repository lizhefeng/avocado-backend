package albums

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/pkg/errors"

	"github.com/lizhefeng/avocado-backend/protocol"
	s "github.com/lizhefeng/avocado-backend/sql"
)

const (
	// ProtocolID is the ID of protocol
	ProtocolID = "albums"
)

// Protocol defines the protocol of indexing blocks
type Protocol struct {
	Store s.Store
}

// NewProtocol creates a new protocol
func NewProtocol(store s.Store) *Protocol {
	return &Protocol{Store: store}
}

// CreateTables creates tables
func (p *Protocol) CreateTables(ctx context.Context) error {
	// create review table
	if _, err := p.Store.GetDB().Exec(fmt.Sprintf("CREATE TABLE IF NOT EXISTS %s "+
		"([review_ID] INTEGER PRIMARY KEY AUTOINCREMENT, [user_ID] BIGINT NOT NULL, [art_ID] BIGINT NOT NULL, [text] TEXT NOT NULL, "+
		"[timestamp] BIGINT NOT NULL, [upvotes] BIGINT NOT NULL)", protocol.ReviewTableName)); err != nil {
		return err
	}
	return nil
}

// PutReview stores a review in the review table
func (p *Protocol) PutReview(ctx context.Context, reviewInput *protocol.ReviewInput) (*protocol.Review, error) {
	insertQuery := fmt.Sprintf("INSERT INTO %s (user_ID, art_ID, [text], timestamp, upvotes) "+
		"VALUES (?, ?, ?, ?, ?)", protocol.ReviewTableName)

	if err := p.Store.Transact(func(tx *sql.Tx) error {
		if _, err := tx.Exec(insertQuery, reviewInput.UserID, reviewInput.ArtID, reviewInput.Text,
			reviewInput.TimeStamp, reviewInput.Upvotes); err != nil {
			return err
		}
		return nil
	}); err != nil {
		return nil, err
	}

	getQuery := fmt.Sprintf("SELECT MAX(review_ID) FROM %s", protocol.ReviewTableName)

	stmt, err := p.Store.GetDB().Prepare(getQuery)
	if err != nil {
		return nil, errors.Wrap(err, "failed to prepare get query")
	}

	var reviewID uint64
	if err = stmt.QueryRow().Scan(&reviewID); err != nil {
		return nil, errors.Wrap(err, "failed to execute get query")
	}
	return &protocol.Review{
		ReviewID:  reviewID,
		UserID:    reviewInput.UserID,
		ArtID:     reviewInput.ArtID,
		Text:      reviewInput.Text,
		TimeStamp: reviewInput.TimeStamp,
		Upvotes:   reviewInput.Upvotes,
	}, nil
}

// GetReview fetches a review from the review table
func (p *Protocol) GetReview(reviewID uint64) (*protocol.Review, error) {
	db := p.Store.GetDB()

	getQuery := fmt.Sprintf("SELECT * FROM %s WHERE review_ID=?", protocol.ReviewTableName)
	stmt, err := db.Prepare(getQuery)
	if err != nil {
		return nil, errors.Wrap(err, "failed to prepare get query")
	}

	rows, err := stmt.Query(reviewID)
	if err != nil {
		return nil, errors.Wrap(err, "failed to execute get query")
	}

	var review protocol.Review
	parsedRows, err := s.ParseSQLRows(rows, &review)
	if err != nil {
		return nil, errors.Wrap(err, "failed to parse results")
	}

	if len(parsedRows) == 0 {
		return nil, protocol.ErrNotExist
	}

	if len(parsedRows) > 1 {
		return nil, errors.New("only one row is expected")
	}

	r := parsedRows[0].(*protocol.Review)
	return r, nil
}
