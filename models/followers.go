package models

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
)

const FollowersCollections = "followers"

type FollowerModel struct {
	ID         primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	UserId     string             `bson:"user_id" json:"user_id"`         //yourid
	FollowerID string             `bson:"follower_id" json:"follower_id"` //twitteruserid
}
