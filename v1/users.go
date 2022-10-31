package route

import (
	"fmt"
	"twitter/controllers"

	"github.com/labstack/echo/v4"
)

func UsersGroup(e *echo.Group) {
	fmt.Println("groups------>")
	e.POST("/signup", controllers.SignUp)
	e.GET("/verify/:otp", controllers.Verify)
	e.POST("/set-password", controllers.SetPassword)
	e.PUT("/upload-profile-pic", controllers.UploadProfilePic)
	e.POST("/set-username", controllers.SetUserName)
	e.POST("/login", controllers.Login)
	//e.POST("/follow", controllers.Follow)
}
