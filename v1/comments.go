package route

import (
	"twitter/controllers"
	"twitter/middleware"

	"github.com/labstack/echo/v4"
)

func CommentsGroup(e *echo.Group) {
	e.POST("/comment/:tweetid", controllers.Comment, middleware.ValidateCustomerToken)
	e.GET("/like-comment/:tweetid", controllers.LikeComment, middleware.ValidateCustomerToken)
	e.GET("/uncomment/:commentid", controllers.Uncomment, middleware.ValidateCustomerToken)
}
