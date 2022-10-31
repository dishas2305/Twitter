package route

import (
	"twitter/controllers"
	"twitter/middleware"

	"github.com/labstack/echo/v4"
)

func TweetsGroup(e *echo.Group) {
	e.POST("/post-tweet/:user_id", controllers.PostTweet, middleware.ValidateCustomerToken)
	e.GET("/my-tweets/:user_id", controllers.GetTweets, middleware.ValidateCustomerToken)
	e.GET("/like-tweet/:tweet_id", controllers.LikeTweet, middleware.ValidateCustomerToken)
	e.GET("/delete-tweet/:tweet_id", controllers.DeleteTweet, middleware.ValidateCustomerToken)
}
