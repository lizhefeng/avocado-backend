// Code generated by github.com/99designs/gqlgen, DO NOT EDIT.

package graphql

type Review struct {
	ReviewID  int    `json:"reviewID"`
	UserID    int    `json:"userID"`
	ArtID     int    `json:"artID"`
	Text      string `json:"text"`
	Timestamp int    `json:"timestamp"`
	Upvotes   int    `json:"upvotes"`
}

type ReviewInput struct {
	UserID    int    `json:"userID"`
	ArtID     int    `json:"artID"`
	Text      string `json:"text"`
	Timestamp int    `json:"timestamp"`
	Upvotes   int    `json:"upvotes"`
}
