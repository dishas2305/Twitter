package route

import (
	"twitter/controllers"
	"twitter/middleware"

	"github.com/labstack/echo/v4"
)

func HomeGroup(e *echo.Group) {
	e.GET("", controllers.Home, middleware.ValidateCustomerToken)
}
