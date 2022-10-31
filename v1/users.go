package route

import (
	"twitter/controllers"
	"twitter/middleware"

	"github.com/labstack/echo/v4"
)

func UsersGroup(e *echo.Group) {
	e.POST("/signup", controllers.SignUp)
	e.GET("/verify/:otp", controllers.Verify)
	e.POST("/set-password", controllers.SetPassword)
	e.PUT("/upload-profile-pic", controllers.UploadProfilePic)
	e.POST("/set-username", controllers.SetUserName)
	e.POST("/login", controllers.Login)
	e.POST("/follow/:follower_id", controllers.Follow, middleware.ValidateCustomerToken)
	e.POST("/unfollow/:follower_id", controllers.Unfollow, middleware.ValidateCustomerToken)
	e.GET("my-followers/:twitter_user_id", controllers.MyFollowers)
}
