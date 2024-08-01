package main

import (
  "context"
  "log"
  "errors"
  "fmt"
	"time"
	"path/filepath"
	"sync"
  "bytes"
  "os"

  "github.com/aws/aws-sdk-go-v2/aws"
  "github.com/aws/aws-sdk-go-v2/config"
  "github.com/aws/aws-sdk-go-v2/service/s3"
  "github.com/spf13/viper"
)

type S3ClientConfig struct {
	RoleArn         string
	RoleSessionName string
	Region          string
	Profile         string
}

var uploadFiles = make(chan string)
var uploadWg = new(sync.WaitGroup)
var uploadBucket = "neuroscan-files"
var directoryIgnore = "/Users/inghamemerson/Code/intralab/neuroscan/applications/neuroscan/backend/public/"
var uploadDirectory = "/Users/inghamemerson/Code/intralab/neuroscan/applications/neuroscan/backend/public/files/neuroscan"
var fileIgnore = []string{".DS_Store", "Thumbs.db", "desktop.ini"}

func main() {
  processUploadDirectory(uploadDirectory)
}

func InitClient(ctx context.Context, bucket string) (*s3.Client, error) {
	clientConfig, err := GetClientConfig(bucket)

	if err != nil {
		log.Fatalf("unable to get client config, %v", err)
	}

	client := CreateS3Client(ctx, clientConfig)

	return client, nil
}

func CreateS3Client(ctx context.Context, clientConfig S3ClientConfig) *s3.Client {
	region := clientConfig.Region
	profile := clientConfig.Profile

	cfg, err := config.LoadDefaultConfig(ctx, config.WithRegion(region), config.WithSharedConfigProfile(profile))

	if err != nil {
		log.Fatalf("unable to load shared profile for sts client config, %v", err)
	}

	Client := s3.NewFromConfig(cfg)

	return Client
}

func GetClientConfig(bucket string) (S3ClientConfig, error) {
	var config S3ClientConfig
	viper.SetConfigName("config")
	viper.SetConfigType("toml")
	viper.AddConfigPath("/etc/vsa/")
	viper.AddConfigPath("~/code/vsa_cli/")
	viper.AddConfigPath(".")
	err := viper.ReadInConfig() // Find and read the config file
	if err != nil {             // Handle errors reading the config file
		return config, errors.New("unable to read config file")
	}

	viperKey := "bucket." + bucket

	roleArn := viper.GetString(viperKey + ".ROLE_ARN")
	roleSessionName := viper.GetString(viperKey + ".ROLE_SESSION_NAME")
	region := viper.GetString(viperKey + ".REGION")
	profile := viper.GetString(viperKey + ".PROFILE")

	config.RoleArn = roleArn
	config.RoleSessionName = roleSessionName
	config.Region = region
	config.Profile = profile

	return config, nil
}

func processUploadDirectory(dir string) {
	startTime := time.Now()
	totalFiles := 0

	for w := 1; w <= 100; w++ {
		go fileUploadWorker()
	}

	defer close(uploadFiles)

	// walk the directory and push each file to the upload channel
	err := filepath.WalkDir(dir, func(path string, d os.DirEntry, err error) error {
		if err != nil {
			log.Println("Unable to walk path: ", err)
			return err
		}

		// if the path is a directory, skip it
		if d.IsDir() {
			return nil
		}

    // if the file is a type we ignore, skip it
    for _, ignore := range fileIgnore {
      if ignore == d.Name() {
        return nil
      }
    }

		// increment the total files
		totalFiles++

		// send the path to the worker pool
		uploadWg.Add(1)

		uploadFiles <- path

		return nil
	})

	if err != nil {
		log.Fatal("Unable to walk directory: ", err)
	}

	// wait for all jobs to finish
	uploadWg.Wait()

	// log the total time it took to process the files
	log.Println("Total time: ", time.Since(startTime))
}

func fileUploadWorker() {
	for file := range uploadFiles {
			processUploadFiles(file)
	}
}

func processUploadFiles(file string) {
	ctx := context.Background()
	client, err := InitClient(ctx, uploadBucket)

	if err != nil {
		log.Fatal("Unable to initialize S3 client: ", err)
	}

	var fileContents []byte


  if _, err := os.Stat(file); os.IsNotExist(err) {
    log.Fatal("File does not exist: ", file)
  }

  fileContents, err = os.ReadFile(file)

  if err != nil {
    log.Fatal("Unable to read file: ", err)
    uploadWg.Done()
    return
  }

  // the key is just the path without the directoryIgnore
  key := file[len(directoryIgnore):]

  fmt.Println("Uploading file: ", key)

  // if the key is empty, skip it
  if key == "" {
    log.Println("Skipping empty key")
    uploadWg.Done()
    return
  }

  _, err = client.PutObject(ctx, &s3.PutObjectInput{
    Bucket: aws.String(uploadBucket),
    Key: aws.String(key),
    Body: bytes.NewReader(fileContents),
  })

  if err != nil {
    log.Fatal("Unable to upload file: ", err)
  }

  uploadWg.Done()
}
