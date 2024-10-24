package services

import (
	"context"
	"sync"

	"github.com/aws/aws-sdk-go-v2/aws"
	v4 "github.com/aws/aws-sdk-go-v2/aws/signer/v4"
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
func (s *StorageService) GetPresignedUrl(key string) (*v4.PresignedHTTPRequest, error) {
	return s.psClient.PresignGetObject(
		context.TODO(),
		&s3.GetObjectInput{
			Bucket: TRACKS_BUCKET_NAME,
			Key:    aws.String(key),
		},
	)
}

// DownloadFile invokes GetObject command with a range if provided
func (s *StorageService) DownloadFile(key string, _range *string) (*s3.GetObjectOutput, error) {
	input := &s3.GetObjectInput{
		Bucket: TRACKS_BUCKET_NAME,
		Key:    aws.String(key),
	}
	if _range != nil {
		input.Range = aws.String(*_range)
	}
	return s.client.GetObject(
		context.TODO(),
		input,
	)
}
