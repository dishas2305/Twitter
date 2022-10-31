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

func Comment(c echo.Context) error {
	comment := &types.CommentBody{}
	_, err := services.GetTweetByID(comment.TweetID)
	if err != nil {
		logger.Error("func_CommentTweet: Record found:", err)
		return utils.HttpErrorResponse(c, utils.GetStatusCode(config.ErrTweetDoesNotExist), config.ErrTweetDoesNotExist)
	}
	err = services.Comment(comment)
	if err != nil {
		logger.Error("func_CommentTweet: ", err)
		return utils.HttpErrorResponse(c, http.StatusBadRequest, err)
	}
	return utils.HttpSuccessResponse(c, http.StatusOK, config.MsgCommentPosted)

}

// func LikeComment(c echo.Context) error {

// }
