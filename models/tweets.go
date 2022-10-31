package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

const TweetsCollections = "tweets"

type TweetModel struct {
	ID       primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	UserID   string             `bson:"user_id,omitempty" json:"user_id"`
	HashTag  string             `bson:"hashtag" json:"hashtag"`
	Tweet    string             `bson:"tweet,omitempty" json:"tweet,omitempty"`
	Picture  string             `bson:"picture,omitempty" json:"picture,omitempty"`
	Likes    int32              `bson:"likes" json:"likes"`
	PostedAt time.Time          `bson:"posted_at" json:"posted_at"`
	Comments string             `bson:"comments,omitempty" json:"comments,omitempty"`
}
