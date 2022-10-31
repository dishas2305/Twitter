package utils

import (
	"mime/multipart"
	"os"
	"strings"
	"time"

	"twitter/config"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
	logger "github.com/sirupsen/logrus"
)

func S3Upload(filename string, file multipart.File) (string, error) {
	if envName := os.Getenv("ENV"); envName != config.Dev && envName != config.Qa && envName != config.Prod {
		logger.Error("not allowed to upload")
		return "", nil
	}

	logger.Info("Upload file to S3")
	uploader, err := config.GetS3Client()
	if err != nil {
		logger.Error("S3Upload: error in GetS3Session is: ", err)
		return "", err
	}

	//upload to the s3 bucket
	uploadInput := &s3manager.UploadInput{
		Bucket: aws.String(os.Getenv("AWS_S3_BUCKET_NAME")),
		ACL:    aws.String("public-read"),
		Key:    aws.String(filename),
		Body:   file,
	}

	if strings.LastIndex(filename, ".pkpass") != -1 {
		expiry := time.Now().UTC().Add(time.Hour * 24)
		uploadInput.Expires = &expiry
	}

	s3Res, err := uploader.Upload(uploadInput)
	if err != nil {
		logger.Error("S3Upload: file name is: ", filename, " error in Upload is: ", err)
		return "", err
	}
	logger.Info("S3Upload: response is: ", s3Res.Location)

	return s3Res.Location, nil
}

func S3DeletedObject(uri string) error {
	if envName := os.Getenv("ENV"); envName != config.Dev && envName != config.Qa && envName != config.Prod {
		logger.Error("not allowed to upload")
		return nil
	}

	logger.Info("Deleting file from S3")
	svc, err := config.GetS3ClientForDeleteObject()
	if err != nil {
		logger.Error("S3DeletedObject: error in get s3 session is: ", err)
		return err
	}

	_, err = svc.DeleteObject(&s3.DeleteObjectInput{
		Bucket: aws.String(os.Getenv("AWS_S3_BUCKET_NAME")),
		Key:    aws.String(uri),
	})
	if err != nil {
		logger.Error("S3DeletedObject: URI is: ", uri, " error in delete object is: ", err)
		return err
	}

	err = svc.WaitUntilObjectNotExists(&s3.HeadObjectInput{
		Bucket: aws.String(os.Getenv("AWS_S3_BUCKET_NAME")),
		Key:    aws.String(uri),
	})
	if err != nil {
		logger.Error("S3DeletedObject: URI is: ", uri, " error in wait until object not exists is: ", err)
		return err
	}

	logger.Info("S3DeletedObject: file deleted")
	return nil
}

func S3DownloadFile(fileName []string, dir string) error {
	if envName := os.Getenv("ENV"); envName != config.Dev && envName != config.Qa && envName != config.Prod {
		logger.Error("not allowed to download")
		return nil
	}

	logger.Info("Downloading file from S3")
	svc, err := config.GetS3ClientForDownloadingFile()
	if err != nil {
		logger.Error("S3DownloadFile: error in get s3 session is: ", err)
		return err
	}

	for _, v := range fileName {
		file, err := os.Create(v)
		if err != nil {
			logger.Error("S3DownloadFile: error in get s3 session is: ", err)
			return err
		}
		defer file.Close()
		_, err = svc.Download(file,
			&s3.GetObjectInput{
				Bucket: aws.String(os.Getenv("AWS_S3_BUCKET_NAME")),
				Key:    aws.String(dir + v),
			})

		if err != nil {
			logger.Error("Unable to download item ", v, err)
			return err
		}
		logger.Info("S3DownloadFile: file downloaded: ", file.Name(), "time: ", time.Now())
	}
	return nil
}
