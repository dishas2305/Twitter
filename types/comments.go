package types

import "go.mongodb.org/mongo-driver/bson/primitive"

type CommentBody struct {
	TweetID primitive.ObjectID `json:"id" example:"1563246835486"`
	Comment string             `json:"comment" example:"comment goes in here"`
}
