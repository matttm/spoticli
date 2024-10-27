package services

import (
	"context"
	"fmt"
	"io"
	"sync"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/s3"
)

// A StorageService interacts directly with s3
type StorageService struct {
	client   *s3.Client
	psClient *s3.PresignClient
}

var storageLock = &sync.Mutex{}
var TRACKS_BUCKET_NAME = aws.String("spoticli-tracks")

var storageService *StorageService

func GetStorageService() *StorageService {
	if storageService == nil {
		storageLock.Lock()
		defer storageLock.Unlock()
		if storageService == nil {
			storageService = &StorageService{}
			storageService.client = s3.NewFromConfig(GetConfigService().CloudConfig)
			storageService.psClient = s3.NewPresignClient(storageService.client)
			println("StorageService Instantiated")
		}
	}
	return storageService
}

// GetPresignedUrl invokes presigned GetObject cmd
func (s *StorageService) GetPresignedUrl(key string) (string, error) {
	res, err := s.psClient.PresignGetObject(
		context.TODO(),
		&s3.GetObjectInput{
			Bucket: TRACKS_BUCKET_NAME,
			Key:    aws.String(key),
		},
	)
	if err != nil {
		panic(err)
	}
	return res.URL, nil
}

// DownloadFile invokes GetObject command with a range if provided
func (s *StorageService) DownloadFile(key string, start, end *int64) ([]byte, error) {
	requestedFrames := make(chan []byte, 1)
	if !isItemCached(key) {
		input := &s3.GetObjectInput{
			Bucket: TRACKS_BUCKET_NAME,
			Key:    aws.String(key),
		}
		res, err := s.client.GetObject(
			context.TODO(),
			input,
		)
		body, err := io.ReadAll(res.Body)
		defer res.Body.Close()
		if err != nil {
			panic(err)
		}
		// the following  blobk is in testing TODO:
		body = ReadID3v2Header(body)
		frames := PartitionMp3Frames(body)
		fmt.Printf("Frame count: %d\n", len(frames))
		// end test NOTE:
		// TODO : put in a goroutine
		cacheItem(key, frames, *start, *end, requestedFrames)
		return getSegmentFromCache(key, *start), nil
	} else {
		return getSegmentFromCache(key, *start), nil
	}
	x := <-requestedFrames
	return x, nil
}
