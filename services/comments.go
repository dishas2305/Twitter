package services

import (
	"context"

	"twitter/models"
	"twitter/storage"
	"twitter/types"

	logger "github.com/sirupsen/logrus"
)

func Comment(c *types.CommentBody) error {
	cm := models.CommentModel{}
	//var response types.VerifyResponse
	cm.ID = c.TweetID
	cm.Comments = c.Comment
	mdb := storage.MONGO_DB
	_, err := mdb.Collection(models.UsersCollections).InsertOne(context.TODO(), cm)
	if err != nil {
		logger.Error("func_SignUp: ", err)
		return err
	}
	return nil
}
