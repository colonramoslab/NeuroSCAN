package main

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"slices"
	"sync"
	"time"

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
var directoryIgnore = "/Users/inghamemerson/Sync/to_upload/"
var uploadDirectory = "/Users/inghamemerson/Sync/to_upload/files/neuroscan"
var bucketFolder string
var validExtensions = []string{".gltf"}

func main() {
	// bucket folder is a string representing the datetime of the upload
	//bucketFolder = time.Now().Format("2006_01_02T15_04_05")
	bucketFolder = "2024_08_21T14_00_28"
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
	var clientConfig S3ClientConfig
	viper.SetConfigName("config")
	viper.SetConfigType("toml")
	viper.AddConfigPath("/etc/neuroscan/")
	viper.AddConfigPath(".")
	err := viper.ReadInConfig() // Find and read the config file
	if err != nil {             // Handle errors reading the config file
		return clientConfig, errors.New("unable to read config file")
	}

	viperKey := "bucket." + bucket

	roleArn := viper.GetString(viperKey + ".ROLE_ARN")
	roleSessionName := viper.GetString(viperKey + ".ROLE_SESSION_NAME")
	region := viper.GetString(viperKey + ".REGION")
	profile := viper.GetString(viperKey + ".PROFILE")

	clientConfig.RoleArn = roleArn
	clientConfig.RoleSessionName = roleSessionName
	clientConfig.Region = region
	clientConfig.Profile = profile

	return clientConfig, nil
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

		// if the extension is not in the valid extensions, skip it
		ext := filepath.Ext(path)

		if !slices.Contains(validExtensions, ext) {
			return nil
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
		uploadWg.Done()
		log.Fatal("Unable to read file: ", err)
		return
	}

	// the key is just the path without the directoryIgnore
	key := file[len(directoryIgnore):]

	// now we add the bucketFolder as a prefix to the key
	key = bucketFolder + "/" + key

	fmt.Println("Uploading file: ", key)

	// if the key is empty, skip it
	if key == "" {
		log.Println("Skipping empty key")
		uploadWg.Done()
		return
	}

	_, err = client.PutObject(ctx, &s3.PutObjectInput{
		Bucket: aws.String(uploadBucket),
		Key:    aws.String(key),
		Body:   bytes.NewReader(fileContents),
	})

	if err != nil {
		log.Fatal("Unable to upload file: ", err)
	}

	uploadWg.Done()
}
