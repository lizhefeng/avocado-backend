package graphql

import (
	"context"
	"github.com/lizhefeng/avocado-backend/keystore"
	"github.com/lizhefeng/avocado-backend/protocol"
	"github.com/lizhefeng/avocado-backend/protocol/albums"
	"github.com/pkg/errors"
) // THIS CODE IS A STARTING POINT ONLY. IT WILL NOT BE UPDATED WITH SCHEMA CHANGES.

// Resolver is the resolver that handles graphql request
type Resolver struct {
	KeyStore *keystore.KeyStore
}

// Mutation returns a query resolver
func (r *Resolver) Mutation() MutationResolver {
	return &mutationResolver{r}
}

// Query returns a query resolver
func (r *Resolver) Query() QueryResolver {
	return &queryResolver{r}
}

type mutationResolver struct{ *Resolver }

// CreateReview creates a review
func (r *mutationResolver) CreateReview(ctx context.Context, review *ReviewInput) (*Review, error) {
	reviewInput := &protocol.ReviewInput{
		UserID:    uint64(review.UserID),
		ArtID:     uint64(review.ArtID),
		Text:      review.Text,
		TimeStamp: uint64(review.Timestamp),
		Upvotes:   uint64(review.Upvotes),
	}
	rev, err := r.KeyStore.PutReview(ctx, albums.ProtocolID, reviewInput)
	if err != nil {
		return nil, errors.Wrap(err, "failed to store review in the underlying db")
	}
	return &Review{
		ReviewID:  int(rev.ReviewID),
		UserID:    int(rev.UserID),
		ArtID:     int(rev.ArtID),
		Text:      rev.Text,
		Timestamp: int(rev.TimeStamp),
		Upvotes:   int(rev.Upvotes),
	}, nil
}

type queryResolver struct{ *Resolver }

// Review returns a review
func (r *queryResolver) Review(ctx context.Context, id int) (*Review, error) {
	rev, err := r.KeyStore.GetReview(albums.ProtocolID, uint64(id))
	if err != nil {
		return nil, errors.Wrap(err, "failed to load review from the underlying db")
	}
	return &Review{
		ReviewID:  int(rev.ReviewID),
		UserID:    int(rev.UserID),
		ArtID:     int(rev.ArtID),
		Text:      rev.Text,
		Timestamp: int(rev.TimeStamp),
		Upvotes:   int(rev.Upvotes),
	}, nil
}
