package controllers

import (
	"fmt"
	"net/http"
	"os"
	"twitter/config"
	"twitter/services"
	"twitter/types"
	"twitter/utils"

	"github.com/labstack/echo/v4"
	logger "github.com/sirupsen/logrus"
)

func SignUp(c echo.Context) error {
	fmt.Println("inside controllers")

	user := &types.SignUpBody{}

	if err := c.Bind(user); err != nil {
		logger.Error("SignUp: Error in binding. Error: ", err)
		return utils.HttpErrorResponse(c, http.StatusBadRequest, config.ErrWrongPayload)
	}

	if err := utils.ValidateStruct(user); err != nil {
		logger.Error("func_SignUp: Error in validating request. Error: ", err)
		return utils.HttpErrorResponse(c, http.StatusBadRequest, err)
	}

	validateMobNum := utils.CheckForNumbers(user.Phone)
	if !validateMobNum {
		logger.Error("func_SignUp: Error :", config.ErrInvalidMobNum)
		return utils.HttpErrorResponse(c, http.StatusBadRequest, config.ErrInvalidMobNum)
	}

	_, err := services.GetUserByMobileNumber(user.Phone)
	if err == nil {
		logger.Error("func_CreateUser: Record found:", err)
		return utils.HttpErrorResponse(c, utils.GetStatusCode(config.ErrDuplicateCustomer), config.ErrDuplicateCustomer)
	}

	// _, err = services.GetUserByEmail(user.Email)
	// if err == nil {
	// 	logger.Error("func_CreateUser: Record found:", err)
	// 	return utils.HttpErrorResponse(c, utils.GetStatusCode(config.ErrDuplicateCustomer), config.ErrDuplicateCustomer)
	// }

	_, err = services.SignUp(user)
	if err != nil {
		logger.Error("func_CreateUser: Error in creating user:", err)
		return utils.HttpErrorResponse(c, http.StatusBadRequest, err)
	}

	return utils.HttpSuccessResponse(c, http.StatusOK, config.MsgUserCreated)

}

func Verify(c echo.Context) error {
	otp := c.Param("otp")
	phone := c.Request().Header.Get("phone")
	fmt.Println(otp)
	verify, err := services.Verify(otp, phone)
	if !verify {
		return utils.HttpErrorResponse(c, utils.GetStatusCode(err), config.ErrInvalidOTP)
	}
	return utils.HttpSuccessResponse(c, http.StatusOK, config.MsgUserVerified)

}

func SetPassword(c echo.Context) error {
	fmt.Println("inside cobntrllers")
	phone := c.Request().Header.Get("Phone")
	result, err := services.GetUserByMobileNumber(phone)
	if err != nil {
		logger.Error("func_SetPassword: Record found:", err)
		return utils.HttpErrorResponse(c, utils.GetStatusCode(config.ErrUserDoesNotExist), config.ErrUserDoesNotExist)
	}
	if result.IsVerified {
		var user types.SetPasswordBody
		err = services.SetPassword(phone, user.Password)
		if err != nil {
			logger.Error("Error in setting password: ", err)
			return utils.HttpErrorResponse(c, utils.GetStatusCode(err), err)
		}
		return utils.HttpSuccessResponse(c, http.StatusOK, config.MsgPasswordSet)

	} else {
		return utils.HttpErrorResponse(c, utils.GetStatusCode(err), config.ErrInvalidOTP)

	}

}

func SetUserName(c echo.Context) error {
	phone := c.Request().Header.Get("Phone")
	handle := c.Request().Header.Get("Handle")
	fmt.Println("handle--------->", handle)
	result, err := services.GetUserByMobileNumber(phone)
	if err != nil {
		logger.Error("func_SetUsername: Record found:", err)
		return utils.HttpErrorResponse(c, utils.GetStatusCode(config.ErrUserDoesNotExist), config.ErrUserDoesNotExist)
	}
	if result.IsVerified {
		err = services.SetUserName(phone, handle)
		if err != nil {
			logger.Error("Error in setting password: ", err)
			return utils.HttpErrorResponse(c, utils.GetStatusCode(err), err)
		}
		return utils.HttpSuccessResponse(c, http.StatusOK, config.MsgUserNameSet)

	} else {
		return utils.HttpErrorResponse(c, utils.GetStatusCode(err), config.ErrInvalidOTP)

	}

}

func UploadProfilePic(c echo.Context) error {
	// Read form fields
	phone := c.FormValue("Phone")
	if len(phone) == 0 {
		return utils.HttpErrorResponse(c, utils.GetStatusCode(config.ErrPhoneKeyNotFound), config.ErrPhoneKeyNotFound)
	}
	// Multipart form
	form, err := c.MultipartForm()
	if err != nil {
		logger.Error("validateFile: Error in MultipartForm. Error: ", err)
		// return nil, err
	}

	// Get data from customer table
	user, err := services.GetUserByMobileNumber(phone)
	if err != nil {
		logger.Error("func_UploadProfilePic: Error in get customer by email. Error: ", err)
		return utils.HttpErrorResponse(c, utils.GetStatusCode(err), err)
	}
	if !user.IsVerified {
		logger.Error("func_UploadProfilePic: customer not verified. Error: ", err)
		return utils.HttpErrorResponse(c, utils.GetStatusCode(config.ErrUserNotVerified), config.ErrUserNotVerified)
	}

	// Read files
	files := form.File["image"]
	if len(files) == 0 {
		logger.Error("func_UploadProfilePic: file not found")
		return utils.HttpErrorResponse(c, utils.GetStatusCode(config.ErrFileNotFound), config.ErrFileNotFound)
	}
	// Validate file (Size, dim etc)
	vI := &utils.ValidateImage{
		Files:           files,
		File:            files[0],
		FileCount:       1,
		EnvFileSizeInKB: os.Getenv("PROFILE_PIC_SIZE_IN_KB"),
		EnvWidth:        os.Getenv("PROFILE_PIC_MIN_WIDTH"),
		EnvHeigth:       os.Getenv("PROFILE_PIC_MIN_HEIGHT"),
	}
	// file, err := utils.ValidateImageFile(files, 1, os.Getenv("PROFILE_PIC_SIZE_IN_KB"), os.Getenv("PROFILE_PIC_MIN_WIDTH"), os.Getenv("PROFILE_PIC_MIN_HEIGHT"))
	file, imgWidth, imgHeight, err := vI.ValidateImageFile()
	if err != nil {
		logger.Error("func_UploadProfilePic: Error in validate image file. Error: ", err)
		return utils.HttpErrorResponse(c, utils.GetStatusCode(err), err)
	}

	// Uploading file to S3

	//utils.S3Upload()
	if err := services.UploadProfilePic(&user, file, imgWidth, imgHeight); err != nil {
		logger.Error("func_UploadProfilePic: Error in upload profile. Error: ", err)
		return utils.HttpErrorResponse(c, utils.GetStatusCode(err), err)
	}

	// Final response
	return utils.HttpSuccessResponse(c, http.StatusOK, "")
}

func Login(c echo.Context) error {
	body := &types.LoginBody{}
	if err := c.Bind(body); err != nil {
		logger.Error("Login: Error in binding. Error: ", err)
		return utils.HttpErrorResponse(c, http.StatusBadRequest, config.ErrWrongPayload)
	}

	if err := utils.ValidateStruct(body); err != nil {
		logger.Error("Login: Error in validating request. Error: ", err)
		return utils.HttpErrorResponse(c, http.StatusBadRequest, err)
	}

	result, err := services.Login(*body)
	if err != nil {
		logger.Error("Login: Error in login. Error: ", err)
		return utils.HttpErrorResponse(c, utils.GetStatusCode(err), err)
	}

	return utils.HttpSuccessResponse(c, http.StatusOK, result)
}

// func Follow(c echo.Context) error {
// 	followerId := c.Request().Header.Get("ID")
// 	_, err := services.GetUserByID(followerId)
// 	if err != nil {
// 		logger.Error("func_Follow: Record found:", err)
// 		return utils.HttpErrorResponse(c, utils.GetStatusCode(config.ErrUserDoesNotExist), config.ErrUserDoesNotExist)
// 	}
// 	var followBody types.FollowBody
// 	_, err = services.GetUserByID(followBody.ID)
// 	if err != nil {
// 		logger.Error("func_follow: Record found:", err)
// 		return utils.HttpErrorResponse(c, utils.GetStatusCode(config.ErrUserDoesNotExist), config.ErrUserDoesNotExist)
// 	}
// 	services.Follow(followBody.ID)

// }
