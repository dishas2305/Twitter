package types

import "go.mongodb.org/mongo-driver/bson/primitive"

type TweetBody struct {
	UserId  string `json:"user_id" example:"45651232789845"`
	HashTag string `json:"hash_tag" example:"#ABCD"`
	Tweet   string `json:"tweet" example:"Tweets on Twitter"`
	Picture string `json:"picture" example:"picId"`
}

type MyTweetsBody struct {
	ID primitive.ObjectID `json:"id" example:"1354687354354"`
}

type LikeTweetBody struct {
	TweetID primitive.ObjectID `json:"id" example:"1354687354354"`
}
