package services

import (
	"context"
	"fmt"

	"twitter/models"
	"twitter/storage"
	"twitter/types"

	logger "github.com/sirupsen/logrus"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func Comment(tweet, user string) error {
	var CommentBody types.CommentBody
	cm := models.CommentModel{}
	cm.Comments = CommentBody.Comment
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
	commentid, err := primitive.ObjectIDFromHex(id)
	_, err = mdb.Collection("comments").UpdateOne(context.TODO(), bson.M{
		"_id": commentid,
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
	collection := mdb.Collection("comments")
	idPrimitive, err := primitive.ObjectIDFromHex(id)

	deleteResult, err := collection.DeleteOne(context.TODO(), bson.M{"_id": idPrimitive})
	if deleteResult.DeletedCount == 0 {
		logger.Error("Error on deleting twwet", err)
		return err
	}
	return nil
}
