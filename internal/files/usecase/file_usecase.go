package usecase

import (
	"context"
	"fmt"
	"mime/multipart"
	"os"
	"ps-gogo-manajer/internal/files/dto"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/aws/aws-sdk-go-v2/service/s3/types"
	"github.com/google/uuid"
	"github.com/pkg/errors"
)

type FileUsecase struct {
	S3Client *s3.Client
}

const (
	JPEG = "image/jpeg"
	JPG  = "image/jpg"
	PNG  = "image/png"
)

var (
	AWS_S3_REGION      = os.Getenv("S3_REGION")
	AWS_S3_ID          = os.Getenv("S3_ID")
	AWS_S3_SECRET_KEY  = os.Getenv("S3_SECRET_KEY")
	AWS_S3_BUCKET_NAME = os.Getenv("S3_BUCKET_NAME")
	nameType           = map[string]string{
		JPEG: ".jpeg",
		JPG:  ".jpg",
		PNG:  ".png",
	}
)

func NewFileUseCase(client *s3.Client) *FileUsecase {
	return &FileUsecase{
		S3Client: client,
	}
}

func (c *FileUsecase) UploadFile(file multipart.File, fileType string) (*dto.FileUploadResponse, error) {
	var response dto.FileUploadResponse
	defer file.Close()

	filename := c.generateFilename(fileType)
	_, err := c.S3Client.PutObject(context.TODO(), &s3.PutObjectInput{
		Bucket: aws.String(AWS_S3_BUCKET_NAME),
		Key:    aws.String(filename),
		ACL:    types.ObjectCannedACLPublicRead,
		Body:   file,
	})

	if err != nil {
		return nil, errors.Wrap(err, "failed to upload file")
	}

	response.FileUrl = c.generateFileUrl(filename)
	return &response, nil
}

func (c *FileUsecase) generateFilename(fileType string) string {
	postfix := nameType[fileType]
	return uuid.New().String() + postfix
}

func (c *FileUsecase) generateFileUrl(filename string) string {
	return fmt.Sprintf(
		"https://%s.s3.%s.amazonaws.com/%s",
		AWS_S3_BUCKET_NAME,
		AWS_S3_REGION,
		filename,
	)
}
