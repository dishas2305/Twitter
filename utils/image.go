package utils

import (
	"errors"
	"image"
	"image/jpeg"
	"io"
	"log"
	"mime/multipart"
	"os"
	"strconv"

	"twitter/config"

	"github.com/nfnt/resize"
	logger "github.com/sirupsen/logrus"
)

type ImageHelp struct {
	SrcFile io.Reader
}

type ValidateImage struct {
	Files           []*multipart.FileHeader
	File            *multipart.FileHeader
	FileCount       int
	EnvFileSizeInKB string
	EnvWidth        string
	EnvHeigth       string
}

type ResizedProfilePic struct {
	ImageData     image.Image
	Width, Height uint
	ImageLoc      string
	AwsFileName   string
}

type ResizedQR struct {
	ImageData     image.Image
	Width, Height uint
	ImageLoc      string
	AwsFileName   string
}

func init() {
	if _, err := os.Stat("images"); os.IsNotExist(err) {
		os.Mkdir("images", 0755)
	}
}

func (ih *ImageHelp) Decode() (image.Image, string, error) {
	imageData, name, err := image.Decode(ih.SrcFile)
	if err != nil {
		logger.Error("func_Decode: Error in image decode: ", err)
		return imageData, name, err
	}

	return imageData, name, nil
}

func GetImageDimension(src multipart.File) (int, int, error) {
	// decode and get config of images
	image, _, err := image.DecodeConfig(src)
	if err != nil {
		logger.Error("GetImageDimension: Error in image DecodeConfig: ", err)
		return 0, 0, err
	}

	return image.Width, image.Height, nil
}

func (vif *ValidateImage) IsValidFileDim() (int, int, error) {
	// Checking file dimension
	src, err := vif.File.Open()
	if err != nil {
		logger.Error("validateFile: Error in file Open: ", err)
		return 0, 0, nil
	}
	defer src.Close()
	imgWidth, imgHeight, err := GetImageDimension(src)
	if err != nil {
		logger.Error("validateFile: Error in image DecodeConfig: ", err)
		return 0, 0, nil
	}
	profileEnvWidth, err := strconv.ParseInt(vif.EnvWidth, 10, 64)
	if err != nil {
		logger.Error("validateFile: Error in profile pic width: ", err)
		return 0, 0, nil
	}
	profileEnvHeight, err := strconv.ParseInt(vif.EnvHeigth, 10, 64)
	if err != nil {
		logger.Error("validateFile: Error in profile pic height: ", err)
		return 0, 0, nil
	}

	logger.Info("profileEnvWidth: ", profileEnvWidth, "profileEnvHeight: ", profileEnvHeight)
	// if int64(imgWidth) < profileEnvWidth || int64(imgHeight) < profileEnvHeight {
	// 	logger.Error("validateFile: Error in profile pic dim check")
	// 	return 0, 0, errors.New(config.ErrProfilePicDim.Error() + vif.EnvWidth + "X" + vif.EnvHeigth)
	// }

	return imgWidth, imgHeight, nil
}

func (vif *ValidateImage) IsValidFileSize() error {
	fileSize, err := strconv.ParseInt(vif.EnvFileSizeInKB, 10, 64)
	if err != nil {
		logger.Error("validateFile: Error in PROFILE_PIC_SIZE_IN_KB env: ", err)
		return err
	}
	if vif.File.Size/1000 > fileSize {
		logger.Error("validateFile: Error in profile file size is more than configured profile file size env")
		return errors.New(config.ErrProfilePicSize.Error() + vif.EnvFileSizeInKB + " KB")
	}
	return nil
}

func (vif *ValidateImage) IsValidFileCount() error {
	// check file count
	if len(vif.Files) == 0 {
		logger.Error("IsFileCountCorrect: file not found")
		return config.ErrFileNotFound
	}
	if len(vif.Files) != vif.FileCount {
		logger.Error("IsFileCountCorrect: can upload single file")
		return config.ErrInvalidFileCount
	}
	return nil
}

func (vif *ValidateImage) ValidateImageFile() (*multipart.FileHeader, int, int, error) {
	// check file count
	if err := vif.IsValidFileCount(); err != nil {
		log.Println("Error in IsValidFileCount")
		log.Println(err)
		logger.Error("func_ValidateImageFile: Error in IsValidFileCount: ", err)
		return nil, 0, 0, err
	}

	// Single pic can upload
	file := vif.Files[0]

	// Checking file size
	if err := vif.IsValidFileSize(); err != nil {
		log.Println("Error in IsValidFileSize")
		log.Println(err)
		logger.Error("func_ValidateImageFile: Error in IsValidFileSize: ", err)
		return nil, 0, 0, err
	}

	// Checking file dimension
	imgWidth, imgHeight, err := vif.IsValidFileDim()
	if err != nil {
		log.Println("Error in is_valid file dim")
		log.Println(err)
		logger.Error("func_ValidateImageFile: Error in is_valid file dim: ", err)
		return nil, 0, 0, err
	}

	return file, imgWidth, imgHeight, nil
}

func (rpp *ResizedProfilePic) UploadResizedProfilePic() error {
	resizedImg := resize.Resize(uint(rpp.Width), uint(rpp.Height), rpp.ImageData, resize.Lanczos3)

	// Creating new image file
	out, err := os.Create(rpp.ImageLoc)
	if err != nil {
		logger.Error("func_UploadResizedProfilePic: Error in os create. Error: ", err)
		return err
	}
	defer out.Close()
	// write new image to file
	err = jpeg.Encode(out, resizedImg, nil)
	if err != nil {
		logger.Error("func_UploadResizedProfilePic: Error in jpeg encode. Error: ", err)
		return err
	}

	// Opening new image file
	//newSrcOrgFile, err := os.Open(rpp.ImageLoc)
	if err != nil {
		logger.Error("func_UploadResizedProfilePic: Error in file open: ", err)
		return err
	}
	//S3Upload(rpp.AwsFileName, newSrcOrgFile)

	if os.Remove(rpp.ImageLoc); err != nil {
		logger.Error("func_UploadResizedProfilePic: Error in os remove: ", err)
		return err
	}
	return nil
}

func (rqr *ResizedQR) UploadResizedQR() error {
	resizedImg := resize.Resize(uint(rqr.Width), uint(rqr.Height), rqr.ImageData, resize.Lanczos3)

	// Creating new image file
	out, err := os.Create(rqr.ImageLoc)
	if err != nil {
		logger.Error("func_UploadResizedQR: Error in os create. Error: ", err)
		return err
	}
	defer out.Close()
	// write new image to file
	err = jpeg.Encode(out, resizedImg, nil)
	if err != nil {
		logger.Error("func_UploadResizedQR: Error in jpeg encode. Error: ", err)
		return err
	}

	// Opening new image file
	//newSrcOrgFile, err := os.Open(rqr.ImageLoc)
	if err != nil {
		logger.Error("func_UploadResizedQR: Error in file open: ", err)
		return err
	}
	//go S3Upload(rqr.AwsFileName, newSrcOrgFile)

	if os.Remove(rqr.ImageLoc); err != nil {
		logger.Error("func_UploadResizedQR: Error in os remove: ", err)
		return err
	}
	return nil
}
