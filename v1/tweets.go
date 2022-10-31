package route

import (
	"twitter/controllers"

	"github.com/labstack/echo/v4"
)

func TweetsGroup(e *echo.Group) {
	e.POST("/post-tweet/:user_id", controllers.PostTweet)
	e.GET("/my-tweets/:user_id", controllers.GetTweets)
	e.GET("/like-tweet/:tweet_id", controllers.LikeTweet)
}
