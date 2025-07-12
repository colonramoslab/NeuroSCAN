package storage

import (
	"bytes"
	"fmt"
	"io"
	"os"

	"neuroscan/pkg/logging"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/joho/godotenv"
)

type Storage struct {
	Client *s3.S3
}

func NewStorage() (*Storage, error) {
	s3Client, err := newClient()
	if err != nil {
		return nil, fmt.Errorf("failed to create S3 client: %w", err)
	}

	return &Storage{
		Client: s3Client,
	}, nil
}

func (s *Storage) PutFile(bucket string, key string, data []byte) error {
	logger := logging.NewLoggerFromEnv()

	body := bytes.NewReader(data)

	_, err := s.Client.PutObject(&s3.PutObjectInput{
		Bucket: aws.String(bucket),
		Key:    aws.String(key),
		Body:   body,
		ACL:    aws.String("public-read"),
	})
	if err != nil {
		logger.Error().Err(err).Msg("Failed to upload file to S3")
		return err
	}

	logger.Info().Msg("File uploaded to S3 successfully")
	return nil
}

func (s *Storage) GetFile(bucket string, key string) ([]byte, error) {
	logger := logging.NewLoggerFromEnv()

	out, err := s.Client.GetObject(&s3.GetObjectInput{
		Bucket: aws.String(bucket),
		Key:    aws.String(key),
	})
	if err != nil {
		logger.Error().Err(err).Msg("Failed to get file from S3")
		return nil, err
	}
	defer out.Body.Close()

	data, err := io.ReadAll(out.Body)
	if err != nil {
		logger.Error().Err(err).Msg("Failed to read file from S3")
	}

	return data, nil
}

func (s *Storage) DeleteFile(bucket string, key string) error {
	logger := logging.NewLoggerFromEnv()

	_, err := s.Client.DeleteObject(&s3.DeleteObjectInput{
		Bucket: aws.String(bucket),
		Key:    aws.String(key),
	})
	if err != nil {
		logger.Error().Err(err).Msg("Failed to delete file from S3")
		return err
	}

	logger.Info().Msg("File deleted from S3 successfully")
	return nil
}

func newClient() (*s3.S3, error) {
	logger := logging.NewLoggerFromEnv()
	err := godotenv.Load()
	if err != nil {
		logger.Info().Err(err).Msg("🤯 failed to load environment variables")
	}

	key := os.Getenv("S3_ACCESS_KEY")
	secret := os.Getenv("S3_ACCESS_SECRET")
	endpoint := os.Getenv("S3_ENDPOINT_URL")
	region := os.Getenv("S3_REGION")

	s3Config := &aws.Config{
		Credentials:      credentials.NewStaticCredentials(key, secret, ""),
		Endpoint:         aws.String(endpoint),
		Region:           aws.String(region),
		S3ForcePathStyle: aws.Bool(false),
	}

	newSession, err := session.NewSession(s3Config)
	if err != nil {
		logger.Fatal().Err(err).Msg("🤯 failed to create AWS session")
		return nil, err
	}

	s3Client := s3.New(newSession)

	return s3Client, nil
}
