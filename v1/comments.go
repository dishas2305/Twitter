package route

import (
	"twitter/controllers"

	"github.com/labstack/echo/v4"
)

func CommentsGroup(e *echo.Group) {
	e.POST("/comment", controllers.Comment)

}
