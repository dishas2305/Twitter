package utils

import (
	crypt "crypto/rand"
	"errors"
	"fmt"
	"math/big"
	"math/rand"

	"os"
	"regexp"
	"strconv"
	"strings"
	"unicode"

	"twitter/config"

	logger "github.com/sirupsen/logrus"
	"golang.org/x/crypto/bcrypt"
)

func FormatInt32(n int32) string {
	return strconv.FormatInt(int64(n), 10)
}

func StringToNumberBig(input string) (int64, error) {
	bi := big.NewInt(0)
	if _, ok := bi.SetString(input, 10); ok {
		return bi.Int64(), nil
	} else {
		return 0, config.ErrNotConvertNumber
	}
}

func CheckForNumbers(str string) bool {
	numCheck := regexp.MustCompile(`^[0-9]+$`)
	return numCheck.MatchString(str)
}

func PasswordValid(mpin string) (bool, error) {
	var (
		letterPresent      bool
		numberPresent      bool
		specialCharPresent bool
		passLen            int
		errorString        string
	)
	MPinLength, err := StringToNumber(os.Getenv("MPIN_LENGTH"))
	if err != nil {
		logger.Error("IsPasswordValid: Parse password length. Error: ", err)
		return false, err
	}

	for _, ch := range mpin {
		switch {
		case unicode.IsNumber(ch):
			numberPresent = true
			passLen++
		case unicode.IsLetter(ch):
			letterPresent = true
			passLen++
		case unicode.IsPunct(ch) || unicode.IsSymbol(ch):
			specialCharPresent = true
			passLen++
		case ch == ' ':
			passLen++
		}
	}
	appendError := func(err string) {
		if len(strings.TrimSpace(errorString)) != 0 {
			errorString += ", " + err
		} else {
			errorString = err
		}
	}
	if letterPresent {
		appendError("MPin cannot conatin character")
		logger.Error("MPin cannot conatin character")
	}
	if !numberPresent {
		appendError("numeric character required")
		logger.Error("numeric character required")
	}
	if specialCharPresent {
		appendError("MPin cannot conatin special character")
		logger.Error("MPin cannot conatin special character")
	}
	if len(mpin) != MPinLength {
		appendError(fmt.Sprintf("MPin length must be between %d digits long", MPinLength))
	}
	if len(errorString) != 0 {
		return false, errors.New(errorString)
	}
	return true, nil
}
func StringToNumber(key string) (int, error) {
	nkey, _ := strconv.Atoi(key)
	return nkey, nil
}

func NumberToString(key int) string {
	skey := strconv.Itoa(key)
	return skey
}

func CheckMPin(mpin, rempin string) int {
	comp := strings.Compare(mpin, rempin)
	return comp
}

func GenOTP(phone string) (string, string, error) {
	max := 9999
	min := 1111
	otp := rand.Intn(max-min) + min
	fmt.Println("otp==============================>", otp)
	strOTP := NumberToString(otp)
	encOTP, err := Encrypt(strOTP, os.Getenv("OTP_ENC_KEY"))
	if err != nil {
		logger.Error("func_CreateUser: Error in encrypt password: ", err)
		return encOTP, strOTP, err
	}
	return encOTP, strOTP, err
}

func CheckMpinMatch(hashedMpin, Mpin string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hashedMpin), []byte(Mpin))
	return err == nil
}

func GetRandomNumber(n uint) (int64, error) {
	i := uint(1)
	nine := "9"
	one := "1"
	for i < n {
		nine = nine + "9"
		one = one + "0"
		i += 1
	}
	nine64, err := StringToNumberBig(nine)
	if err != nil {
		logger.Error("GetRandomNumber: Parse nine Error: ", err)
		return 0, err
	}
	one64, err := StringToNumberBig(one)
	if err != nil {
		logger.Error("GetRandomNumber: Parse one Error: ", err)
		return 0, err
	}

	nBig, err := crypt.Int(crypt.Reader, big.NewInt(nine64-one64))
	if err != nil {
		logger.Error("GetRandomNumber: Error: ", err)
		return 0, err
	}
	return nBig.Int64() + one64, nil
}

func GetS3Folder(n int) string {
	i := 0
	result := ""
	splitedString := strings.Split(strconv.Itoa(n), "")
	for i < len(splitedString) {
		result += "/" + splitedString[i]
		i += 1
	}
	return result
}

func GetBucketKey(folderName string, customerId int, imgView string, imgName string) string {
	if customerId == 0 {
		return "/" + folderName + "/" + imgView + "/" + imgName
	} else {
		return "/" + folderName + GetS3Folder(customerId) + "/" + imgView + "/" + imgName
	}
}

func GetS3URL(folderName string, customerId int, imgView string, imgName string) string {
	awsHostURL := "https://" + os.Getenv("AWS_S3_BUCKET_NAME") + ".s3." + os.Getenv("AWS_S3_REGION") + ".amazonaws.com"
	awsFileKey := GetBucketKey(folderName, customerId, imgView, imgName)
	return awsHostURL + awsFileKey
}
