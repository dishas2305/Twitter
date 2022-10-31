package config

import "errors"

const (
	MsgUserCreated     = "User Created. OTP sent for Verification"
	MsgLoginSuccessful = "Logged in successfully"
	MsgOTPSent         = "OTP sent to Mobile Number"
	MsgMpinChanged     = "MPin changed successfully"
	MsgSiteAdded       = "Site added successfully"
	MsgSiteUpdated     = "Site details updated sucessfully"
	MsgUserVerified    = "User verified sucessfully"
	MsgPasswordSet     = "Password set sucessfully"
	MsgUserNameSet     = "UserName set sucessfully"
	MsgPostTweet       = "Post tweeted sucessfully"
	MsgTweetLiked      = "Tweet liked"
	MsgCommentPosted   = "Comment posted sucessfully"
)

var (
	ErrMissingBasicAuth            = errors.New("Authorization must be required in header")
	ErrWrongPayload                = errors.New("Wrong payload, please try again")
	ErrRecordNotFound              = errors.New("Record not found")
	ErrParameterMissing            = errors.New("Parameter missing")
	ErrTokenMissing                = errors.New("Error token missing")
	ErrInvalidHashKey              = errors.New("Invalid hash key")
	ErrInvalidHttpMethod           = errors.New("Invalid http method")
	ErrHttpCallBadRequest          = errors.New("Bad request")
	ErrHttpCallUnauthorized        = errors.New("Unauthorized")
	ErrHttpCallNotFound            = errors.New("Call not found")
	ErrHttpCallInternalServerError = errors.New("Internal server error")
	ErrWentWrong                   = errors.New("Something went wrong")
	ErrInvalidMobNum               = errors.New("Invalid mobile number")
	ErrInvalidPasswordFormat       = errors.New("Invalid password format")
	ErrDuplicateCustomer           = errors.New("User already exists with this mobile number")
	ErrVerKeyNotFound              = errors.New("verify key not found")
	ErrEmailAlreadyVerified        = errors.New("Email already verified")
	ErrUserDoesNotExist            = errors.New("User does not exist with this credentials")
	ErrEmailNotVerified            = errors.New(" Email not verified")
	ErrInvalidToken                = errors.New("Invalid token")
	ErrRefOnly                     = errors.New("Reference Only")
	ErrInvalidPassword             = errors.New("Invalid Password")
	ErrMPinDoNotMatch              = errors.New("MPins do not match")
	ErrInvalidOTP                  = errors.New("Invalid OTP")
	ErrDuplicateSite               = errors.New("Site with this URL already exists in this folder")
	ErrSiteNotFound                = errors.New("Site not found")
	ErrURLKeyNotFound              = errors.New("Site URL not found")
	ErrFileNotFound                = errors.New("File not found")
	ErrPhoneNotFound               = errors.New("Phone not found")
	ErrPhoneKeyNotFound            = errors.New("Phone Key not found")
	ErrUserNotVerified             = errors.New("User not verified")
	ErrProfilePicSize              = errors.New("Error in profile picture size")
	ErrInvalidFileCount            = errors.New("Error in file count")
	ErrNotConvertNumber            = errors.New("Cannot convert number")
	ErrTweetDoesNotExist           = errors.New("Tweet does not exist")
)
