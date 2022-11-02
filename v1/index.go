package route

import (
	"github.com/labstack/echo/v4"
)

func InitializeRoutes(e *echo.Group) {
	gHome := e.Group("/home")
	HomeGroup(gHome)

	gUsers := e.Group("/users")
	UsersGroup(gUsers)

	gTweets := e.Group("/tweets")
	TweetsGroup(gTweets)

	gComments := e.Group("/comments")
	CommentsGroup(gComments)
}
