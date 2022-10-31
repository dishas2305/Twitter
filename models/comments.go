package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

const CommentsCollections = "comments"

type CommentModel struct {
	ID       primitive.ObjectID `bson:"_id,omitempty" json:"id"`

	Likes    int                `bson:"likes" json:"likes"`
	PostedAt time.Time          `bson:"posted_at" json:"posted_at"`
	Comments string             `bson:"comments,omitempty" json:"comments,omitempty"`
}
