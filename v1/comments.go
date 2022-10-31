package route

import (
	"twitter/controllers"

	"github.com/labstack/echo/v4"
)

func CommentsGroup(e *echo.Group) {
	e.POST("/comment/:tweetid", controllers.Comment)
	e.GET("/like-comment/:tweetid", controllers.LikeComment)
	e.GET("/uncomment/:commentid", controllers.Uncomment)

}
