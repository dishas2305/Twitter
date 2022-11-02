package services

import (
	"context"
	"twitter/models"
	"twitter/storage"

	logger "github.com/sirupsen/logrus"
	"go.mongodb.org/mongo-driver/bson"
)

func Home() ([]models.TweetModel, error) {
	var tweets []models.TweetModel
	mdb := storage.MONGO_DB
	filter := bson.M{}

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
