package storage

import (
	"context"
	"io"
	"os"
	"path/filepath"
	"strconv"
	"sync"
	"time"

	"neuroscan/internal/toolshed"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	// "github.com/aws/aws-sdk-go-v2/service/s3/types"
	"github.com/rs/zerolog/log"
	"github.com/schollz/progressbar/v3"
)

func DownloadBucketFile(bucket string, destination string, key string, prefix string, workers int) error {
	// if key and prefix are both set, throw an error
	if key != "" && prefix != "" {
		log.Fatal().Msg("Cannot set both key and prefix")
	}

	// if neither key nor prefix are set, throw an error
	if key == "" && prefix == "" {
		log.Fatal().Msg("Must set either key or prefix")
	}

	// if key is set, download the file
	if key != "" {
		log.Info().Msg("Downloading file: " + key + " from bucket: " + bucket)
		processDownloadBucketFile(bucket, destination, key)
		return nil
	}

	// if prefix is set, download all files with that prefix
	if prefix != "" {
		log.Info().Msg("Downloading files with prefix: " + prefix + " from bucket: " + bucket)
		processDownloadBucketPrefix(bucket, destination, prefix, workers)
		return nil
	}

	return nil
}

func processDownloadBucketPrefix(bucket string, destination string, prefix string, workers int) {
	ctx := context.Background()
	startTime := time.Now()

	client, err := InitClient(ctx)

	if err != nil {
		log.Fatal().Err(err).Msg("unable to setup client")
	}

	// ensure destination ends with a slash
	if destination[len(destination)-1:] != "/" {
		destination = destination + "/"
	}

	log.Info().Msg("determining files to download")

	jobs := ListBucketItems(ctx, client, bucket, prefix)

	log.Info().Msg("downloading " + strconv.Itoa(len(jobs)) + " files")

	filesChannel := make(chan string, len(jobs))

	var wg sync.WaitGroup

	bar := progressbar.Default(int64(len(jobs)), "Downloading")

	totalItems := int64(0)

	for range workers {

		go func() {

			for key := range filesChannel {
				processDownloadBucketFile(bucket, destination, key)
				wg.Done()
				bar.Add(1)
				totalItems++
			}
		}()
	}

	for _, item := range jobs {
		wg.Add(1)
		filesChannel <- item
	}

	close(filesChannel)

	wg.Wait()

	log.Info().Msg("Downloaded " + strconv.FormatInt(totalItems, 10) + " files in: " + time.Since(startTime).String())
}

func processDownloadBucketFile(bucket string, destination string, key string) error {
	ctx := context.Background()

	client, err := InitClient(ctx)

	if err != nil {
		log.Fatal().Err(err).Msg("unable to setup client")
		return err
	}

	// ensure destination ends with a slash
	if destination[len(destination)-1:] != "/" {
		destination = destination + "/"
	}

	file, err := client.GetObject(ctx, &s3.GetObjectInput{
		Bucket: aws.String(bucket),
		Key:    aws.String(key),
	})

	if err != nil {
		log.Fatal().Err(err).Msg("unable to get object")
		return err
	}

	defer file.Body.Close()

	toolshed.CreateDirectory(destination+filepath.Dir(key), 0755)
	newFile, err := os.Create(destination + key)

	if err != nil {
		log.Fatal().Err(err).Msg("unable to create file")
		return err
	}

	defer newFile.Close()

	f, err := os.OpenFile(destination+key, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)

	if err != nil {
		log.Fatal().Err(err).Msg("unable to open file")
		return err
	}

	defer f.Close()

	body, err := io.ReadAll(file.Body)

	if err != nil {
		log.Fatal().Err(err).Msg("unable to read body")
		return err
	}

	_, err = f.Write(body)

	if err != nil {
		log.Fatal().Err(err).Msg("unable to write to file")
		return err
	}

	return nil
}
