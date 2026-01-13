package config

import (
	"context"
	"log"
	"os"

	AWSConfig "github.com/aws/aws-sdk-go-v2/config"
	AWSCredentials "github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/service/s3"
)

type S3Client *s3.Client

var (
	AWS_S3_REGION      = os.Getenv("S3_REGION")
	AWS_S3_ID          = os.Getenv("S3_ID")
	AWS_S3_SECRET_KEY  = os.Getenv("S3_SECRET_KEY")
	AWS_S3_BUCKET_NAME = os.Getenv("S3_BUCKET_NAME")
)

func NewS3Client() *s3.Client {
	config, err := AWSConfig.LoadDefaultConfig(
		context.TODO(),
		AWSConfig.WithRegion(AWS_S3_REGION),
		AWSConfig.WithCredentialsProvider(
			AWSCredentials.NewStaticCredentialsProvider(
				AWS_S3_ID,
				AWS_S3_SECRET_KEY,
				""),
		),
	)
	if err != nil {
		log.Fatal("unable connect to S3 Client", err.Error())
	}

	client := s3.NewFromConfig(config)
	return client
}
