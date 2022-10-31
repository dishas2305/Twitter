package controllers

import (
	"net/http"
	"twitter/config"
	"twitter/services"
	"twitter/types"
	"twitter/utils"

	"github.com/labstack/echo/v4"
	logger "github.com/sirupsen/logrus"
)

func PostTweet(c echo.Context) error {
	userId := c.Param("user_id")
	user := &types.TweetBody{}

	if err := c.Bind(user); err != nil {
		logger.Error("Post Tweet: Error in binding. Error: ", err)
		return utils.HttpErrorResponse(c, http.StatusBadRequest, config.ErrWrongPayload)
	}

	if err := utils.ValidateStruct(user); err != nil {
		logger.Error("Post Tweet: Error in validating request. Error: ", err)
		return utils.HttpErrorResponse(c, http.StatusBadRequest, err)
	}

	_, err := services.PostTweet(user, userId)
	if err != nil {
		logger.Error("Post Tweet: Error in posting tweet:", err)
		return utils.HttpErrorResponse(c, http.StatusBadRequest, err)
	}

	return utils.HttpSuccessResponse(c, http.StatusOK, config.MsgPostTweet)

}

func GetTweets(c echo.Context) error {
	userId := c.Param("user_id")
	//var gettweetbody types.MyTweetsBody
	// _, err := services.GetTweetByID(gettweetbody.ID)
	// if err != nil {
	// 	logger.Error("func_GetTweets: Record found:", err)
	// 	return utils.HttpErrorResponse(c, utils.GetStatusCode(config.ErrTweetDoesNotExist), config.ErrTweetDoesNotExist)
	// }
	result, err := services.GetTweets(userId)
	if err != nil {
		logger.Error("func_GetTweets: Record found:", err)
		return utils.HttpErrorResponse(c, http.StatusBadRequest, err)
	}
	return utils.HttpSuccessResponse(c, http.StatusOK, result)
}

func LikeTweet(c echo.Context) error {
	tweetId := c.Param("tweet_id")
	err := services.LikeTweet(tweetId)
	if err != nil {
		logger.Error("func_LikeTweet: ", err)
		return utils.HttpErrorResponse(c, http.StatusBadRequest, err)
	}
	return utils.HttpSuccessResponse(c, http.StatusOK, config.MsgTweetLiked)
}

func DeleteTweet(c echo.Context) error {
	tweetid := c.Param("tweet_id")
	err := services.DeleteTweet(tweetid)
	if err != nil {
		logger.Error("func_DeleteTweet: ", err)
		return utils.HttpErrorResponse(c, http.StatusBadRequest, err)
	}
	return utils.HttpSuccessResponse(c, http.StatusOK, config.MsgTweetDeleted)
}
