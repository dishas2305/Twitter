package controllers

import (
	"fmt"
	"net/http"
	"twitter/config"
	"twitter/services"
	"twitter/types"
	"twitter/utils"

	"github.com/labstack/echo/v4"
	logger "github.com/sirupsen/logrus"
)

func Comment(c echo.Context) error {
	tweetID := c.Param("tweetid")
	userID := c.Request().Header.Get("userID")
	var comment types.CommentBody
	fmt.Println(comment)
	_, err := services.GetTweetByID(tweetID)
	if err != nil {
		logger.Error("func_CommentTweet: Record found:", err)
		return utils.HttpErrorResponse(c, utils.GetStatusCode(config.ErrTweetDoesNotExist), config.ErrTweetDoesNotExist)
	}
	err = services.Comment(tweetID, userID, comment.Comment)
	if err != nil {
		logger.Error("func_CommentTweet: ", err)
		return utils.HttpErrorResponse(c, http.StatusBadRequest, err)
	}
	return utils.HttpSuccessResponse(c, http.StatusOK, config.MsgCommentPosted)

}

func LikeComment(c echo.Context) error {
	fmt.Println("inside controllers")
	tweetId := c.Param("tweetid")
	_, err := services.GetTweetByID(tweetId)
	if err != nil {
		logger.Error("func_CommentTweet: Record found:", err)
		return utils.HttpErrorResponse(c, utils.GetStatusCode(config.ErrTweetDoesNotExist), config.ErrTweetDoesNotExist)
	}
	err = services.LikeComment(tweetId)
	if err != nil {
		logger.Error("func_LikeComment: ", err)
		return utils.HttpErrorResponse(c, http.StatusBadRequest, err)
	}
	return utils.HttpSuccessResponse(c, http.StatusOK, config.MsgCommentLiked)
}

func Uncomment(c echo.Context) error {
	commentid := c.Param("commentid")
	err := services.Uncomment(commentid)
	if err != nil {
		logger.Error("func_UnComment: ", err)
		return utils.HttpErrorResponse(c, http.StatusBadRequest, err)
	}
	return utils.HttpSuccessResponse(c, http.StatusOK, config.MsgCommentDeleted)
}
