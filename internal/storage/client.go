package storage

import (
	"context"
	"os"

	"neuroscan/internal/logging"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/joho/godotenv"
)

type S3ClientConfig struct {
	Bucket 					string
	RoleArn         string
	RoleSessionName string
	Region          string
	Profile         string
}

func CreateS3Client(ctx context.Context, clientConfig S3ClientConfig) *s3.Client {
	// roleArn := clientConfig.RoleArn
	// roleSessionName := clientConfig.RoleSessionName
	region := clientConfig.Region
	profile := clientConfig.Profile

	logger := logging.FromContext(ctx)

	cfg, err := config.LoadDefaultConfig(ctx, config.WithRegion(region), config.WithSharedConfigProfile(profile))

	if err != nil {
		logger.Fatal().Err(err).Msg("unable to load shared profile for sts client config")
	}

	Client := s3.NewFromConfig(cfg)

	return Client
}

func GetClientConfig(ctx context.Context) (S3ClientConfig, error) {
	var config S3ClientConfig

	err := godotenv.Load()
	logger := logging.FromContext(ctx)

	if err != nil {
		logger.Err(err).Msg("unable to load environment variables")
		return config, err
	}

	bucket := os.Getenv( "S3_BUCKET")
	roleArn := os.Getenv( "S3_ROLE_ARN")
	roleSessionName := os.Getenv( "S3_ROLE_SESSION_NAME")
	region := os.Getenv( "S3_REGION")
	profile := os.Getenv( "S3_PROFILE")

	config.Bucket = bucket
	config.RoleArn = roleArn
	config.RoleSessionName = roleSessionName
	config.Region = region
	config.Profile = profile

	return config, nil
}

func ValidateConfig(config S3ClientConfig) bool {
	if config.Bucket == "" || config.RoleArn == "" || config.RoleSessionName == "" || config.Region == "" || config.Profile == "" {
		return false
	}

	return true
}

func ListBucketItems(ctx context.Context, client *s3.Client, bucket string, prefix string) []string {
	var items []string
	logger := logging.FromContext(ctx)

	input := &s3.ListObjectsV2Input{
		Bucket: aws.String(bucket),
	}

	if prefix != "" {
		input.Prefix = aws.String(prefix)
	}

	paginator := s3.NewListObjectsV2Paginator(client, input)

	for paginator.HasMorePages() {
		page, err := paginator.NextPage(ctx)

		if err != nil {
			logger.Fatal().Err(err).Msg("unable to list objects")
		}

		for _, object := range page.Contents {
			items = append(items, *object.Key)
		}
	}

	return items
}

func InitClient(ctx context.Context) (*s3.Client, error) {
	clientConfig, err := GetClientConfig(ctx)
	logger := logging.FromContext(ctx)

	if err != nil {
		logger.Fatal().Err(err).Msg("unable to get client config")
	}

	if !ValidateConfig(clientConfig) {
		logger.Fatal().Msg("invalid client config")
	}

	client := CreateS3Client(ctx, clientConfig)

	return client, nil
}
