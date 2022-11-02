package controllers

import (
	"net/http"
	"twitter/services"
	"twitter/utils"

	"github.com/labstack/echo/v4"
	logger "github.com/sirupsen/logrus"
)

func Home(c echo.Context) error {
	result, err := services.Home()
	if err != nil {
		logger.Error("Home Page: Error in loading homepage:", err)
		return utils.HttpErrorResponse(c, http.StatusBadRequest, err)
	}
	return utils.HttpSuccessResponse(c, http.StatusOK, result)
}
