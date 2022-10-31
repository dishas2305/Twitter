package services

import (
	"context"
	"fmt"
	"image"
	"mime/multipart"
	"os"
	"strconv"
	"time"

	"twitter/config"
	"twitter/models"
	"twitter/storage"
	"twitter/types"
	"twitter/utils"

	"github.com/dgrijalva/jwt-go"
	logger "github.com/sirupsen/logrus"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func GetUserByMobileNumber(phone string) (models.UserModel, error) {
	var user models.UserModel
	mdb := storage.MONGO_DB

	filter := bson.M{
		"phone": phone,
	}

	result := mdb.Collection(models.UsersCollections).FindOne(context.TODO(), filter)
	err := result.Decode(&user)
	if err != nil {
		logger.Error("func_GetUserByMobileNUmber: Error in ", err)
		return user, err
	}
	return user, nil
}

func GetUserByID(Id string) (models.UserModel, error) {
	var user models.UserModel
	mdb := storage.MONGO_DB

	filter := bson.M{
		"_id": Id,
	}

	result := mdb.Collection(models.UsersCollections).FindOne(context.TODO(), filter)
	err := result.Decode(&user)
	if err == nil {
		logger.Error("func_GetUserByMobileNUmber: Error in ", err)
		return user, err
	}
	return user, nil
}

func GetUserByEmail(email string) (models.UserModel, error) {
	var user models.UserModel
	mdb := storage.MONGO_DB

	filter := bson.M{
		"email": email,
	}

	result := mdb.Collection(models.UsersCollections).FindOne(context.TODO(), filter)
	err := result.Decode(&user)
	if err != nil {
		logger.Error("func_GetUserByEmail: Error in ", err)
		return user, err
	}
	return user, nil
}
func GetUserByUserName(username string) (models.UserModel, error) {
	var user models.UserModel
	mdb := storage.MONGO_DB

	filter := bson.M{
		"user_name": username,
	}

	result := mdb.Collection(models.UsersCollections).FindOne(context.TODO(), filter)
	err := result.Decode(&user)
	if err != nil {
		logger.Error("func_GetUserByUseranme: Error in ", err)
		return user, err
	}
	return user, nil
}

func SignUp(c *types.SignUpBody) (types.SignUpResponse, error) {
	fmt.Println("inside services")
	var response types.SignUpResponse
	um := models.UserModel{}
	//var response types.VerifyResponse
	um.Name = c.Name
	um.Phone = c.Phone
	um.Email = c.Email
	um.DOB = c.DOB
	mdb := storage.MONGO_DB
	_, err := mdb.Collection(models.UsersCollections).InsertOne(context.TODO(), um)
	if err != nil {
		logger.Error("func_SignUp: ", err)
		return response, err
	}
	encOTP, _, err := utils.GenOTP(um.Phone)
	response.OTP = encOTP
	um.OTP = encOTP
	filter := bson.M{
		"phone": um.Phone,
	}
	update := bson.M{"$set": bson.M{"otp": encOTP}}
	_, err = mdb.Collection(models.UsersCollections).UpdateOne(context.TODO(), filter, update)
	if err != nil {
		logger.Error("func_CreateUser: ", err)
		return response, err
	}

	return response, nil
}

func Verify(otp, phone string) (bool, error) {
	var user models.UserModel
	mdb := storage.MONGO_DB
	filter := bson.M{
		"phone": phone,
	}
	result := mdb.Collection(models.UsersCollections).FindOne(context.TODO(), filter)
	err := result.Decode(&user)
	if err != nil {
		logger.Error("func_User not found: Error in ", err)
		return false, err
	}
	fmt.Println("user--------->", user)
	decOTP, err := utils.Decrypt(user.OTP, os.Getenv("OTP_ENC_KEY"))
	fmt.Println("decotp------------------------>", decOTP)
	if err != nil {
		logger.Error("func_encryption error: Error in ", err)
	}

	if decOTP == otp {
		_, err = mdb.Collection("users").UpdateOne(context.TODO(), bson.M{
			"phone": phone,
		}, bson.D{
			{"$set", bson.D{{"is_verified", true}}},
		}, options.Update().SetUpsert(true))
	} else {
		return false, err
	}
	return true, nil
}

func SetPassword(phone, password string) error {
	fmt.Println("insde services")
	mdb := storage.MONGO_DB
	filter := bson.M{
		"phone": phone,
	}
	encPassword, err := utils.Encrypt(password, os.Getenv("PASSWORD_ENC_KEY"))
	update := bson.M{"$set": bson.M{"password": encPassword}}
	_, err = mdb.Collection(models.UsersCollections).UpdateOne(context.TODO(), filter, update)
	if err != nil {
		logger.Error("func_SetPassword: ", err)
		return err
	}
	return err
}

func SetUserName(phone, handle string) error {
	fmt.Println("insde services")
	mdb := storage.MONGO_DB
	filter := bson.M{
		"phone": phone,
	}
	update := bson.M{"$set": bson.M{"handle": handle}}
	_, err := mdb.Collection(models.UsersCollections).UpdateOne(context.TODO(), filter, update)
	if err != nil {
		logger.Error("func_SetPassword: ", err)
		return err
	}
	return err
}

func UploadProfilePic(user *models.UserModel, file *multipart.FileHeader, imgWidth int, imgHeight int) error {
	newExt := "jpg"

	randomName, err := utils.GetRandomNumber(8)
	if err != nil {
		logger.Error("func_UploadProfilePic: Error in get random number: ", err)
		return err
	}

	srcFile, err := file.Open()
	if err != nil {
		logger.Error("func_UploadProfilePic: Error in file open: ", err)
		return err
	}
	defer srcFile.Close()

	// decode
	imageData, _, err := image.Decode(srcFile)
	if err != nil {
		logger.Error("func_UploadProfilePic: Error in image decode: ", err)
		return err
	}

	newFileName := ""
	if user.ProfilePicture == "0" || len(user.ProfilePicture) == 0 {
		newFileName = strconv.FormatInt(randomName, 10) + "." + newExt
	} else {
		newFileName = user.ProfilePicture
	}
	imageLoc := "images" + "/" + newFileName

	userId := fmt.Sprintln(user.ID)
	userIdint, _ := strconv.Atoi(userId)

	thumbnailPath := utils.GetBucketKey("profile", userIdint, "thumbnail", newFileName)
	displayPath := utils.GetBucketKey("profile", userIdint, "display", newFileName)

	fmt.Println("thumbnailPath===>", thumbnailPath, "displayPath==>", displayPath)
	rpp := &utils.ResizedProfilePic{
		ImageData:   imageData,
		Width:       uint(imgWidth),
		Height:      uint(imgHeight),
		ImageLoc:    imageLoc,
		AwsFileName: thumbnailPath,
	}
	if err := rpp.UploadResizedProfilePic(); err != nil {
		logger.Error("func_UploadProfilePic: Error in Upload-Resized-ProfilePic for thumbnail. Error: ", err)
		return err
	}

	rpp.Width = uint(imgWidth)
	rpp.Height = uint(imgHeight)

	rpp.AwsFileName = displayPath
	if err := rpp.UploadResizedProfilePic(); err != nil {
		logger.Error("func_UploadProfilePic: Error in Upload-Resized-ProfilePic for display. Error: ", err)
		return err
	}

	user.ProfilePicture = newFileName

	fmt.Println("newFileName===>>", newFileName)
	if err := UpdateCustomer(user, newFileName); err != nil {
		logger.Error("func_UploadProfilePic: Error in update customer. Error: ", err)
		return err
	}

	return nil
}

func UpdateCustomer(c *models.UserModel, newFilename string) error {
	mdb := storage.MONGO_DB
	filter := bson.M{
		"ID": c.ID,
	}
	update := bson.M{"$set": bson.M{"profile": newFilename}}
	_, err := mdb.Collection(models.UsersCollections).UpdateOne(context.TODO(), filter, update)
	if err != nil {
		logger.Error("func_SetUsername: ", err)
		return err
	}
	return err
}

func Login(payload types.LoginBody) (types.LoginOutput, error) {
	var loginOutput types.LoginOutput
	user, err := GetUserByUserName(payload.UserName)
	if err == nil {
		logger.Error("GetUserByMobileNumber: Error in fetching customer by mobile number. Error: ", err)
		return loginOutput, err
	}
	reqPass := payload.Password
	decPass, err := utils.Decrypt(user.Password, os.Getenv("PASSWORD_ENC_KEY"))
	if reqPass == decPass {
		fmt.Println("password matches")
		token, err := GenerateToken(user)
		fmt.Println("token---------------------->", token)
		if err != nil {
			logger.Error("GenerateToken: Error in generating the token Error: ", err)
			return loginOutput, err
		}
		rtoken, err := GenerateRefreshToken(user)
		fmt.Println("refersh token-------------------------->", rtoken)
		if err != nil {
			logger.Error("GenerateToken: Error in generating the refresh token Error: ", err)
			return loginOutput, err
		}
		mdb := storage.MONGO_DB
		filter := bson.M{
			"_id": user.ID,
		}
		update := bson.M{"$set": bson.M{"refresh_token": rtoken}}
		_, err = mdb.Collection(models.UsersCollections).UpdateOne(context.TODO(), filter, update)
		if err != nil {
			logger.Error("func_Login: ", err)
			return loginOutput, err
		}
		loginOutput.Token = token
		loginOutput.RefreshToken = rtoken
		return loginOutput, nil
	} else {
		return loginOutput, config.ErrInvalidPassword
	}
}

func GenerateToken(userResult models.UserModel) (string, error) {
	token := jwt.New(jwt.SigningMethodHS256)
	tokenExpiredBy, _ := strconv.Atoi(os.Getenv("ACCESS_TOKEN_EXPIRY"))
	claims := token.Claims.(jwt.MapClaims)
	claims["user_id"] = userResult.ID
	claims["user_name"] = userResult.UserName
	claims["exp"] = time.Now().Add(time.Hour * 24 * time.Duration(tokenExpiredBy)).Unix()

	return token.SignedString([]byte(os.Getenv("CUSTOMER_JWT_SECRET_KEY")))
}

func GenerateRefreshToken(userResult models.UserModel) (string, error) {
	token := jwt.New(jwt.SigningMethodHS256)
	tokenExpiredBy, _ := strconv.Atoi(os.Getenv("REFRESH_TOKEN_EXPIRY"))
	claims := token.Claims.(jwt.MapClaims)
	claims["user_id"] = userResult.ID
	claims["is_refresh"] = true
	claims["exp"] = time.Now().Add(time.Hour * 24 * time.Duration(tokenExpiredBy)).Unix()
	refershToken, err := token.SignedString([]byte(os.Getenv("CUSTOMER_JWT_SECRET_KEY")))
	return refershToken, err
}

func Follow(myid, yourid string) error {
	mdb := storage.MONGO_DB
	_, err := mdb.Collection("users").UpdateOne(context.TODO(), bson.M{
		"_id": yourid,
	}, bson.D{
		{"$inc", bson.D{{"followers", 1}}},
	}, options.Update().SetUpsert(true))
	if err != nil {
		logger.Error("func_GetTweets: ", err)
		return err
	}

	_, err = mdb.Collection("users").UpdateOne(context.TODO(), bson.M{
		"_id": myid,
	}, bson.D{
		{"$inc", bson.D{{"following", 1}}},
	}, options.Update().SetUpsert(true))
	if err != nil {
		logger.Error("func_GetTweets: ", err)
		return err
	}

	um := models.FollowerModel{}
	um.UserId = yourid
	um.FollowerID = myid
	_, err = mdb.Collection(models.FollowersCollections).InsertOne(context.TODO(), um)
	if err != nil {
		logger.Error("func_SignUp: ", err)
		return err
	}
	return nil
}

func Unfollow(myid, yourid string) error {
	mdb := storage.MONGO_DB
	_, err := mdb.Collection("users").UpdateOne(context.TODO(), bson.M{
		"_id": yourid,
	}, bson.D{
		{"$dec", bson.D{{"followers", 1}}},
	}, options.Update().SetUpsert(true))
	if err != nil {
		logger.Error("func_Unfollow: ", err)
		return err
	}

	_, err = mdb.Collection("users").UpdateOne(context.TODO(), bson.M{
		"_id": myid,
	}, bson.D{
		{"$dec", bson.D{{"following", 1}}},
	}, options.Update().SetUpsert(true))
	if err != nil {
		logger.Error("func_Unfollow: ", err)
		return err
	}
	return nil
}

func MyFollowers(id string) ([]models.FollowerModel, error) {
	var followers []models.FollowerModel
	mdb := storage.MONGO_DB
	filter := bson.M{
		"user_id": id,
	}
	result, err := mdb.Collection(models.FollowersCollections).Find(context.TODO(), filter)
	if err != nil {
		logger.Error("func_GetFollowers: ", err)
		return followers, err
	}
	if err := result.All(context.Background(), &followers); err != nil {
		logger.Error("func_GetFollowers: error cur.All() step ", err)
		return nil, err
	}
	return followers, nil
}
