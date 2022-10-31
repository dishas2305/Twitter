package services

import (
	"context"
	"fmt"

	"twitter/models"
	"twitter/storage"

	logger "github.com/sirupsen/logrus"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func Comment(tweet, user, comment string) error {
	cm := models.CommentModel{}
	cm.Comments = comment
	cm.TweetID = tweet
	cm.UserID = user
	mdb := storage.MONGO_DB
	_, err := mdb.Collection(models.CommentsCollections).InsertOne(context.TODO(), cm)
	if err != nil {
		logger.Error("func_SignUp: ", err)
		return err
	}
	return nil
}

func LikeComment(id string) error {
	fmt.Println("inside services")
	mdb := storage.MONGO_DB
	_, err := mdb.Collection("comments").UpdateOne(context.TODO(), bson.M{
		"_id": id,
	}, bson.D{
		{"$inc", bson.D{{"likes", 1}}},
	}, options.Update().SetUpsert(true))
	if err != nil {
		logger.Error("func_LikeComment: ", err)
		return err
	}
	return nil
}

func Uncomment(id string) error {
	mdb := storage.MONGO_DB
	_, err := mdb.Collection("commments").DeleteOne(context.TODO(), bson.M{"_id": "id"})
	if err != nil {
		logger.Error("func_GetTweets: ", err)
		return err
	}
	return nil

}
