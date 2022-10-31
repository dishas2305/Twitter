package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

const CommentsCollections = "comments"

type CommentModel struct {
	ID       primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	TweetID  string             `bson:"tweet_id" json:"tweet_id"`
	UserID   string             `bson:"user_id" json:"user_id"`
	Likes    int                `bson:"likes" json:"likes"`
	PostedAt time.Time          `bson:"posted_at" json:"posted_at"`
	Comments string             `bson:"comments" json:"comments"`
}
