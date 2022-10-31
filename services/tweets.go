package services

import (
	"context"
	"twitter/models"
	"twitter/storage"
	"twitter/types"

	logger "github.com/sirupsen/logrus"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func GetTweetByID(Id string) (models.TweetModel, error) {
	var tweet models.TweetModel
	mdb := storage.MONGO_DB
	filter := bson.M{
		"_id": Id,
	}

	result := mdb.Collection(models.TweetsCollections).FindOne(context.TODO(), filter)
	err := result.Decode(&tweet)
	if err == nil {
		logger.Error("func_GetTweetByID: Error in ", err)
		return tweet, err
	}
	return tweet, nil
}

func PostTweet(c *types.TweetBody, userid string) (models.TweetModel, error) {
	tm := models.TweetModel{}
	tm.HashTag = c.HashTag
	tm.Tweet = c.Tweet
	//tm.Picture =
	tm.UserID = userid

	mdb := storage.MONGO_DB
	_, err := mdb.Collection(models.TweetsCollections).InsertOne(context.TODO(), tm)
	if err != nil {
		logger.Error("func_PostTweet: ", err)
		return tm, err
	}
	return tm, nil
}

func GetTweets(id string) ([]models.TweetModel, error) {
	var tweets []models.TweetModel
	mdb := storage.MONGO_DB
	filter := bson.M{
		"user_id": id,
	}
	result, err := mdb.Collection(models.TweetsCollections).Find(context.TODO(), filter)
	if err != nil {
		logger.Error("func_GetTweets: ", err)
		return tweets, err
	}
	if err := result.All(context.Background(), &tweets); err != nil {
		logger.Error("func_GetCities: error cur.All() step ", err)
		return nil, err
	}
	return tweets, nil
}

func LikeTweet(id string) error {
	mdb := storage.MONGO_DB
	_, err := mdb.Collection("tweets").UpdateOne(context.TODO(), bson.M{
		"_id": id,
	}, bson.D{
		{"$inc", bson.D{{"likes", 1}}},
	}, options.Update().SetUpsert(true))
	if err != nil {
		logger.Error("func_GetTweets: ", err)
		return err
	}

	return nil
}

func DeleteTweet(id string) error {
	mdb := storage.MONGO_DB
	collection := mdb.Collection("tweets")
	idPrimitive, err := primitive.ObjectIDFromHex(id)

	deleteResult, err := collection.DeleteOne(context.TODO(), bson.M{"_id": idPrimitive})
	if deleteResult.DeletedCount == 0 {
		logger.Error("Error on deleting twwet", err)
		return err
	}
	return nil
}
